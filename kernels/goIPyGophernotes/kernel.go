package main

import (
	"context"
	"encoding/json"
	"fmt"
  "io"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	"github.com/go-zeromq/zmq4"
	"golang.org/x/xerrors"
)

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
type KernelLanguageInfo struct {
	Name              string `json:"name"`
	Version           string `json:"version"`
	MIMEType          string `json:"mimetype"`
	FileExtension     string `json:"file_extension"`
	PygmentsLexer     string `json:"pygments_lexer"`
	CodeMirrorMode    string `json:"codemirror_mode"`
	NBConvertExporter string `json:"nbconvert_exporter"`
}

// HelpLink stores data to be displayed in the help menu of the notebook.
type HelpLink struct {
	Text string `json:"text"`
	URL  string `json:"url"`
}

// KernelInfo holds information about the igo kernel, for kernel_info_reply messages.
type KernelInfo struct {
	ProtocolVersion       string             `json:"protocol_version"`
	Implementation        string             `json:"implementation"`
	ImplementationVersion string             `json:"implementation_version"`
	LanguageInfo          KernelLanguageInfo `json:"language_info"`
	Banner                string             `json:"banner"`
	HelpLinks             []HelpLink         `json:"help_links"`
}

// shutdownReply encodes a boolean indication of shutdown/restart.
type ShutdownReply struct {
	Restart bool `json:"restart"`
}

const (
	KernelStarting = "starting"
	KernelBusy     = "busy"
	KernelIdle     = "idle"
)

type AdaptorImpl interface {

  // GetKernelInfo returns the KernelInfo for this kernel implementation.
  //
  GetKernelInfo() KernelInfo
  
  // Get the possible completions for the word at cursorPos in the code. 
  //
  GetCodeWordCompletions(code string, cursorPos int) (int, int, []string)

  // Setup the Display callback by recording the msgReceipt information
  // for later use by what ever callback implements the "Display" function. 
  //
  SetupDisplayCallback(receipt msgReceipt)
  
  // Teardown the Display callback by removing the current msgReceipt
  // information and setting things back to what ever default the 
  // implementation uses.
  //
  TeardownDisplayCallback()
  
  // Evaluate (and remove) any implmenation specific special commands BEFORE 
  // the code gets evaluated by the interpreter. The `outErr` variable 
  // contains the stdOut and stdErr which can be used to capture the stdOut 
  // and stdErr of any external commands run by these special commands. 
  //
  EvaluateRemoveSpecialCommands(outErr OutErr, code string) string

  // Evaluate the code and return the results as a Data object.
  //
  EvaluateCode(code string) (rtnData Data, err error)

}

type IPyKernel struct {
  adaptor AdaptorImpl
}

func NewIPyKernel(anAdaptor AdaptorImpl) *IPyKernel {
  return &IPyKernel{ adaptor: anAdaptor }
}

// RunWithSocket invokes the `run` function after acquiring the 
// `Socket.Lock` and releases the lock when done. 
//
func (s *Socket) RunWithSocket(run func(socket zmq4.Socket) error) error {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	return run(s.Socket)
}

// IPyKernel::Run is the main entry point to start the kernel.
//
func (kernel *IPyKernel) Run(connectionFile string) {

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
	sockets, err := PrepareSockets(connInfo)
	if err != nil {
		log.Fatal(err)
	}

  // TODO connect all channel handlers to a WaitGroup to ensure shutdown 
  // before returning from runKernel. 

  // Start up the heartbeat handler.
	StartHeartbeat(sockets.HBSocket, &sync.WaitGroup{})

  // TODO gracefully shutdown the heartbeat handler on kernel shutdown by 
  // closing the chan returned by startHeartbeat. 

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

			kernel.HandleShellMsg(msgReceipt{msg, ids, sockets})

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

			kernel.HandleShellMsg(msgReceipt{msg, ids, sockets})
		}
	}
}

// prepareSockets sets up the ZMQ sockets through which the kernel
// will communicate.
//
func PrepareSockets(connInfo ConnectionInfo) (SocketGroup, error) {
	// Initialize the socket group.
	var (
		sg  SocketGroup
		err error
		ctx = context.Background()
	)

  // Create the shell socket, a request-reply socket that may receive 
  // messages from multiple frontend for code execution, introspection, 
  // auto-completion, etc. 
  //
  sg.ShellSocket.Socket = zmq4.NewRouter(ctx)
	sg.ShellSocket.Lock = &sync.Mutex{}

  // Create the control socket. This socket is a duplicate of the shell 
  // socket where messages on this channel should jump ahead of queued 
  // messages on the shell socket. 
  //
	sg.ControlSocket.Socket = zmq4.NewRouter(ctx)
	sg.ControlSocket.Lock = &sync.Mutex{}

  // Create the stdin socket, a request-reply socket used to request user 
  // input from a front-end. This is analogous to a standard input stream. 
  //
	sg.StdinSocket.Socket = zmq4.NewRouter(ctx)
	sg.StdinSocket.Lock = &sync.Mutex{}

  // Create the iopub socket, a publisher for broadcasting data like 
  // stdout/stderr output, displaying execution results or errors, kernel 
  // status, etc. to connected subscribers. 
  //
	sg.IOPubSocket.Socket = zmq4.NewPub(ctx)
	sg.IOPubSocket.Lock = &sync.Mutex{}

  // Create the heartbeat socket, a request-reply socket that only allows 
  // alternating recv-send (request-reply) calls. It should echo the byte 
  // strings it receives to let the requester know the kernel is still 
  // alive. 
  //
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
func (kernel *IPyKernel) HandleShellMsg(receipt msgReceipt) {
	// Tell the front-end that the kernel is working and when finished notify the
	// front-end that the kernel is idle again.
	if err := receipt.PublishKernelStatus(KernelBusy); err != nil {
		log.Printf("Error publishing kernel status 'busy': %v\n", err)
	}
	defer func() {
		if err := receipt.PublishKernelStatus(KernelIdle); err != nil {
			log.Printf("Error publishing kernel status 'idle': %v\n", err)
		}
	}()

	switch receipt.Msg.Header.MsgType {
	case "kernel_info_request":
		if err := kernel.HandleKernelInfoRequest(receipt); err != nil {
			log.Fatal(err)
		}
	case "complete_request":
		if err := kernel.HandleCompleteRequest(receipt); err != nil {
			log.Fatal(err)
		}
	case "execute_request":
		if err := kernel.HandleExecuteRequest(receipt); err != nil {
			log.Fatal(err)
		}
	case "shutdown_request":
		kernel.HandleShutdownRequest(receipt)
	default:
		log.Println("Unhandled shell message: ", receipt.Msg.Header.MsgType)
	}
}

func (kernel *IPyKernel) HandleKernelInfoRequest(receipt msgReceipt) error {
	return receipt.Reply(
    "kernel_info_reply",
    kernel.adaptor.GetKernelInfo(),
  )
}

func (kernel *IPyKernel) HandleCompleteRequest(receipt msgReceipt) error {
	// Extract the data from the request.
	reqcontent := receipt.Msg.Content.(map[string]interface{})
	code := reqcontent["code"].(string)
	cursorPos := int(reqcontent["cursor_pos"].(float64))

	// autocomplete the code at the cursor position
  cursorStart, cursorEnd, matches := 
    kernel.adaptor.GetCodeWordCompletions(code, cursorPos)
  
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

// handleExecuteRequest runs code from an execute_request method,
// and sends the various reply messages.
//
func (kernel *IPyKernel) HandleExecuteRequest(receipt msgReceipt) error {

	// Extract the data from the request.
	reqcontent := receipt.Msg.Content.(map[string]interface{})
	code := reqcontent["code"].(string)
	silent := reqcontent["silent"].(bool)

	if !silent {
		ExecCounter++
	}

	// Prepare the map that will hold the reply content.
	content := make(map[string]interface{})
	content["execution_count"] = ExecCounter

	// Tell the front-end what the kernel is about to execute.
	if err := receipt.PublishExecutionInput(ExecCounter, code); err != nil {
		log.Printf("Error publishing execution input: %v\n", err)
	}

	// Redirect the standard out from the REPL.
	oldStdout := os.Stdout
	rOut, wOut, err := os.Pipe()
	if err != nil {
		return err
	}
	os.Stdout = wOut

	// Redirect the standard error from the REPL.
	oldStderr := os.Stderr
	rErr, wErr, err := os.Pipe()
	if err != nil {
		return err
	}
	os.Stderr = wErr

	var writersWG sync.WaitGroup
	writersWG.Add(2)

	jupyterStdOut := JupyterStreamWriter{StreamStdout, &receipt}
	jupyterStdErr := JupyterStreamWriter{StreamStderr, &receipt}
	outerr := OutErr{&jupyterStdOut, &jupyterStdErr}

	// Forward all data written to stdout/stderr to the front-end.
	go func() {
		defer writersWG.Done()
		io.Copy(&jupyterStdOut, rOut)
	}()

	go func() {
		defer writersWG.Done()
		io.Copy(&jupyterStdErr, rErr)
	}()

  kernel.adaptor.SetupDisplayCallback(receipt)
  defer kernel.adaptor.TeardownDisplayCallback()
  
  // evaluate and remove any special commands
  code = kernel.adaptor.EvaluateRemoveSpecialCommands(outerr, code)
  
	// eval
	data, executionErr := kernel.adaptor.EvaluateCode(code)

	// Close and restore the streams.
	wOut.Close()
	os.Stdout = oldStdout

	wErr.Close()
	os.Stderr = oldStderr

	// Wait for the writers to finish forwarding the data.
	writersWG.Wait()

	if executionErr == nil {
		// if the only non-nil value should be auto-rendered graphically, render it

		content["status"] = "ok"
		content["user_expressions"] = make(map[string]string)

		if !silent && len(data.Data) != 0 {
			// Publish the result of the execution.
			if err := receipt.PublishExecutionResult(ExecCounter, data); err != nil {
				log.Printf("Error publishing execution result: %v\n", err)
			}
		}
	} else {
		content["status"] = "error"
		content["ename"] = "ERROR"
		content["evalue"] = executionErr.Error()
		content["traceback"] = nil

		if err := receipt.PublishExecutionError(
      executionErr.Error(),
      []string{executionErr.Error()},
    ); err != nil {
			log.Printf("Error publishing execution error: %v\n", err)
		}
	}

	// Send the output back to the notebook.
	return receipt.Reply("execute_reply", content)
}

// handleShutdownRequest sends a "shutdown" message.
//
func (kernel *IPyKernel) HandleShutdownRequest(receipt msgReceipt) {
	content := receipt.Msg.Content.(map[string]interface{})
	restart := content["restart"].(bool)

	reply := ShutdownReply{
		Restart: restart,
	}

	if err := receipt.Reply("shutdown_reply", reply); err != nil {
		log.Fatal(err)
	}

	log.Println("Shutting down in response to shutdown_request")
	os.Exit(0)
}

// startHeartbeat starts a go-routine for handling heartbeat ping messages 
// sent over the given `hbSocket`. The `wg`'s `Done` method is invoked 
// after the thread is completely shutdown. To request a shutdown the 
// returned `shutdown` channel can be closed. 
//
func StartHeartbeat(hbSocket Socket, wg *sync.WaitGroup) (shutdown chan struct{}) {
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
