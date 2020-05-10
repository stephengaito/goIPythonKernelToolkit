package main

import (
	"context"
	"encoding/json"
//	"errors"
	"fmt"
//	"go/ast"
//	"io"
	"io/ioutil"
	"log"
	"os"
//	"os/exec"
//	"reflect"
//	"runtime"
//	"strings"
	"sync"
	"time"

	"github.com/go-zeromq/zmq4"
	"golang.org/x/xerrors"

//	"github.com/cosmos72/gomacro/ast2"
//	"github.com/cosmos72/gomacro/base"
//	basereflect "github.com/cosmos72/gomacro/base/reflect"
//	interp "github.com/cosmos72/gomacro/fast"
//	"github.com/cosmos72/gomacro/xreflect"

	// compile and link files generated in imports/
	_ "github.com/stephengaito/goIPythonKernelToolkit/kernels/goIPyGophernotes/imports"
)

type InterpreterImpl interface {

  // SendKernelInfo sends a kernel_info_reply message.
  //
  SendKernelInfo(receipt msgReceipt) error
  
  // HandleExecuteRequest runs code from an execute_request method,
  // and sends the various reply messages.
  //
  HandleExecuteRequest(receipt msgReceipt) error
  
  // Get the possible completions for the word at cursorPos in the code. 
  //
  GetCodeWordCompletions(code string, cursorPos int) (int, int, []string)
}

// SHOULD MOVE TO ADAPTOR?
//
// ExecCounter is incremented each time we run user code in the notebook.
var ExecCounter int

// ConnectionInfo stores the contents of the kernel connection
// file created by Jupyter.
type ConnectionInfo struct {
	SignatureScheme string `json:"signature_scheme"`
	Transport       string `json:"transport"`
	StdinPort       int    `json:"stdin_port"`
	ControlPort     int    `json:"control_port"`
	IOPubPort       int    `json:"iopub_port"`
	HBPort          int    `json:"hb_port"`
	ShellPort       int    `json:"shell_port"`
	Key             string `json:"key"`
	IP              string `json:"ip"`
}

// Socket wraps a zmq socket with a lock which should be used to control write access.
type Socket struct {
	Socket zmq4.Socket
	Lock   *sync.Mutex
}

// SocketGroup holds the sockets needed to communicate with the kernel,
// and the key for message signing.
type SocketGroup struct {
	ShellSocket   Socket
	ControlSocket Socket
	StdinSocket   Socket
	IOPubSocket   Socket
	HBSocket      Socket
	Key           []byte
}

// KernelLanguageInfo holds information about the language that this kernel executes code in.
type kernelLanguageInfo struct {
	Name              string `json:"name"`
	Version           string `json:"version"`
	MIMEType          string `json:"mimetype"`
	FileExtension     string `json:"file_extension"`
	PygmentsLexer     string `json:"pygments_lexer"`
	CodeMirrorMode    string `json:"codemirror_mode"`
	NBConvertExporter string `json:"nbconvert_exporter"`
}

// HelpLink stores data to be displayed in the help menu of the notebook.
type helpLink struct {
	Text string `json:"text"`
	URL  string `json:"url"`
}

// KernelInfo holds information about the igo kernel, for kernel_info_reply messages.
type kernelInfo struct {
	ProtocolVersion       string             `json:"protocol_version"`
	Implementation        string             `json:"implementation"`
	ImplementationVersion string             `json:"implementation_version"`
	LanguageInfo          kernelLanguageInfo `json:"language_info"`
	Banner                string             `json:"banner"`
	HelpLinks             []helpLink         `json:"help_links"`
}

// shutdownReply encodes a boolean indication of shutdown/restart.
type shutdownReply struct {
	Restart bool `json:"restart"`
}

const (
	kernelStarting = "starting"
	kernelBusy     = "busy"
	kernelIdle     = "idle"
)

// RunWithSocket invokes the `run` function after acquiring the `Socket.Lock` and releases the lock when done.
func (s *Socket) RunWithSocket(run func(socket zmq4.Socket) error) error {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	return run(s.Socket)
}

// runKernel is the main entry point to start the kernel.
func runKernel(kernel InterpreterImpl, connectionFile string) {

	// Parse the connection info.
	var connInfo ConnectionInfo

	connData, err := ioutil.ReadFile(connectionFile)
	if err != nil {
		log.Fatal(err)
	}

	if err = json.Unmarshal(connData, &connInfo); err != nil {
		log.Fatal(err)
	}

	// Set up the ZMQ sockets through which the kernel will communicate.
	sockets, err := prepareSockets(connInfo)
	if err != nil {
		log.Fatal(err)
	}

	// TODO connect all channel handlers to a WaitGroup to ensure shutdown before returning from runKernel.

	// Start up the heartbeat handler.
	startHeartbeat(sockets.HBSocket, &sync.WaitGroup{})

	// TODO gracefully shutdown the heartbeat handler on kernel shutdown by closing the chan returned by startHeartbeat.

	type msgType struct {
		Msg zmq4.Msg
		Err error
	}

	var (
		shell = make(chan msgType)
		stdin = make(chan msgType)
		ctl   = make(chan msgType)
		quit  = make(chan int)
	)

	defer close(quit)
	poll := func(msgs chan msgType, sck zmq4.Socket) {
		defer close(msgs)
		for {
			msg, err := sck.Recv()
			select {
			case msgs <- msgType{Msg: msg, Err: err}:
			case <-quit:
				return
			}
		}
	}

	go poll(shell, sockets.ShellSocket.Socket)
	go poll(stdin, sockets.StdinSocket.Socket)
	go poll(ctl, sockets.ControlSocket.Socket)

	// Start a message receiving loop.
	for {
		select {
		case v := <-shell:
			// Handle shell messages.
			if v.Err != nil {
				log.Println(v.Err)
				continue
			}

			msg, ids, err := WireMsgToComposedMsg(v.Msg.Frames, sockets.Key)
			if err != nil {
				log.Println(err)
				return
			}

			handleShellMsg(msgReceipt{msg, ids, sockets}, kernel)

		case <-stdin:
			// TODO Handle stdin socket.
			continue

		case v := <-ctl:
			if v.Err != nil {
				log.Println(v.Err)
				return
			}

			msg, ids, err := WireMsgToComposedMsg(v.Msg.Frames, sockets.Key)
			if err != nil {
				log.Println(err)
				return
			}

			handleShellMsg(msgReceipt{msg, ids, sockets}, kernel)
		}
	}
}

// prepareSockets sets up the ZMQ sockets through which the kernel
// will communicate.
func prepareSockets(connInfo ConnectionInfo) (SocketGroup, error) {
	// Initialize the socket group.
	var (
		sg  SocketGroup
		err error
		ctx = context.Background()
	)

	// Create the shell socket, a request-reply socket that may receive messages from multiple frontend for
	// code execution, introspection, auto-completion, etc.
	sg.ShellSocket.Socket = zmq4.NewRouter(ctx)
	sg.ShellSocket.Lock = &sync.Mutex{}

	// Create the control socket. This socket is a duplicate of the shell socket where messages on this channel
	// should jump ahead of queued messages on the shell socket.
	sg.ControlSocket.Socket = zmq4.NewRouter(ctx)
	sg.ControlSocket.Lock = &sync.Mutex{}

	// Create the stdin socket, a request-reply socket used to request user input from a front-end. This is analogous
	// to a standard input stream.
	sg.StdinSocket.Socket = zmq4.NewRouter(ctx)
	sg.StdinSocket.Lock = &sync.Mutex{}

	// Create the iopub socket, a publisher for broadcasting data like stdout/stderr output, displaying execution
	// results or errors, kernel status, etc. to connected subscribers.
	sg.IOPubSocket.Socket = zmq4.NewPub(ctx)
	sg.IOPubSocket.Lock = &sync.Mutex{}

	// Create the heartbeat socket, a request-reply socket that only allows alternating recv-send (request-reply)
	// calls. It should echo the byte strings it receives to let the requester know the kernel is still alive.
	sg.HBSocket.Socket = zmq4.NewRep(ctx)
	sg.HBSocket.Lock = &sync.Mutex{}

	// Bind the sockets.
	address := fmt.Sprintf("%v://%v:%%v", connInfo.Transport, connInfo.IP)
	err = sg.ShellSocket.Socket.Listen(fmt.Sprintf(address, connInfo.ShellPort))
	if err != nil {
		return sg, xerrors.Errorf("could not listen on shell-socket: %w", err)
	}

	err = sg.ControlSocket.Socket.Listen(fmt.Sprintf(address, connInfo.ControlPort))
	if err != nil {
		return sg, xerrors.Errorf("could not listen on control-socket: %w", err)
	}

	err = sg.StdinSocket.Socket.Listen(fmt.Sprintf(address, connInfo.StdinPort))
	if err != nil {
		return sg, xerrors.Errorf("could not listen on stdin-socket: %w", err)
	}

	err = sg.IOPubSocket.Socket.Listen(fmt.Sprintf(address, connInfo.IOPubPort))
	if err != nil {
		return sg, xerrors.Errorf("could not listen on iopub-socket: %w", err)
	}

	err = sg.HBSocket.Socket.Listen(fmt.Sprintf(address, connInfo.HBPort))
	if err != nil {
		return sg, xerrors.Errorf("could not listen on hbeat-socket: %w", err)
	}

	// Set the message signing key.
	sg.Key = []byte(connInfo.Key)

	return sg, nil
}

// handleShellMsg responds to a message on the shell ROUTER socket.
func handleShellMsg(receipt msgReceipt, kernel InterpreterImpl) {
	// Tell the front-end that the kernel is working and when finished notify the
	// front-end that the kernel is idle again.
	if err := receipt.PublishKernelStatus(kernelBusy); err != nil {
		log.Printf("Error publishing kernel status 'busy': %v\n", err)
	}
	defer func() {
		if err := receipt.PublishKernelStatus(kernelIdle); err != nil {
			log.Printf("Error publishing kernel status 'idle': %v\n", err)
		}
	}()

	switch receipt.Msg.Header.MsgType {
	case "kernel_info_request":
		if err := kernel.SendKernelInfo(receipt); err != nil {
			log.Fatal(err)
		}
	case "complete_request":
		if err := handleCompleteRequest(kernel, receipt); err != nil {
			log.Fatal(err)
		}
	case "execute_request":
		if err := kernel.HandleExecuteRequest(receipt); err != nil {
			log.Fatal(err)
		}
	case "shutdown_request":
		handleShutdownRequest(receipt)
	default:
		log.Println("Unhandled shell message: ", receipt.Msg.Header.MsgType)
	}
}

func handleCompleteRequest(kernel InterpreterImpl, receipt msgReceipt) error {
	// Extract the data from the request.
	reqcontent := receipt.Msg.Content.(map[string]interface{})
	code := reqcontent["code"].(string)
	cursorPos := int(reqcontent["cursor_pos"].(float64))

	// autocomplete the code at the cursor position
  cursorStart, cursorEnd, matches := 
    kernel.GetCodeWordCompletions(code, cursorPos)
  
	// prepare the reply
	content := make(map[string]interface{})

	if len(matches) == 0 {
		content["ename"] = "ERROR"
		content["evalue"] = "no completions found"
		content["traceback"] = nil
		content["status"] = "error"
	} else {
		content["cursor_start"] = float64(cursorStart)
		content["cursor_end"] = float64(cursorEnd)
		content["matches"] = matches
		content["status"] = "ok"
	}

	return receipt.Reply("complete_reply", content)
}

// handleShutdownRequest sends a "shutdown" message.
func handleShutdownRequest(receipt msgReceipt) {
	content := receipt.Msg.Content.(map[string]interface{})
	restart := content["restart"].(bool)

	reply := shutdownReply{
		Restart: restart,
	}

	if err := receipt.Reply("shutdown_reply", reply); err != nil {
		log.Fatal(err)
	}

	log.Println("Shutting down in response to shutdown_request")
	os.Exit(0)
}

// startHeartbeat starts a go-routine for handling heartbeat ping messages sent over the given `hbSocket`. The `wg`'s
// `Done` method is invoked after the thread is completely shutdown. To request a shutdown the returned `shutdown` channel
// can be closed.
func startHeartbeat(hbSocket Socket, wg *sync.WaitGroup) (shutdown chan struct{}) {
	quit := make(chan struct{})

	// Start the handler that will echo any received messages back to the sender.
	wg.Add(1)
	go func() {
		defer wg.Done()

		type msgType struct {
			Msg zmq4.Msg
			Err error
		}

		msgs := make(chan msgType)

		go func() {
			defer close(msgs)
			for {
				msg, err := hbSocket.Socket.Recv()
				select {
				case msgs <- msgType{msg, err}:
				case <-quit:
					return
				}
			}
		}()

		timeout := time.NewTimer(500 * time.Second)
		defer timeout.Stop()

		for {
			timeout.Reset(500 * time.Second)
			select {
			case <-quit:
				return
			case <-timeout.C:
				continue
			case v := <-msgs:
				hbSocket.RunWithSocket(func(echo zmq4.Socket) error {
					if v.Err != nil {
						log.Fatalf("Error reading heartbeat ping bytes: %v\n", v.Err)
						return v.Err
					}

					// Send the received byte string back to let the front-end know that the kernel is alive.
					if err := echo.Send(v.Msg); err != nil {
						log.Printf("Error sending heartbeat pong bytes: %b\n", err)
						return err
					}

					return nil
				})
			}
		}
	}()

	return quit
}
