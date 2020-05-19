package goIPyGoMacroAdaptor

// (proto)Adaptor 

import(
//	"context"
//	"encoding/json"
	"errors"
	"fmt"
	"go/ast"
	"io"
	"io/ioutil"
	"log"
//	"os"
	"os/exec"
	"reflect"
	"runtime"
	"strings"
	"sync"
//	"time"

  tk "github.com/stephengaito/goIPythonKernelToolkit/goIPyKernel"

	"github.com/cosmos72/gomacro/ast2"
	"github.com/cosmos72/gomacro/base"
	basereflect "github.com/cosmos72/gomacro/base/reflect"
	gomacro "github.com/cosmos72/gomacro/fast"
	"github.com/cosmos72/gomacro/xreflect"

	// compile and link files generated in imports/
	_ "github.com/stephengaito/goIPythonKernelToolkit/kernels/goIPyGophernotes/imports"
)

const (
	// Version defines the goIPyGophernotes version.
	Version string = "1.0.0"
)

type GoAdaptor struct {
	ir      *gomacro.Interp
	display *gomacro.Import
	// map name -> HTMLer, JSONer, Renderer...
	// used to convert interpreted types to one of these interfaces
	render map[string]xreflect.Type
}

func NewGoAdaptor() *GoAdaptor {
	// Create a new interpreter for evaluating notebook code.
	ir := gomacro.New()

	// Throw out the error/warning messages that gomacro outputs writes to these streams.
	ir.Comp.Stdout = ioutil.Discard
	ir.Comp.Stderr = ioutil.Discard

	// Inject the "display" package to render HTML, JSON, PNG, JPEG, SVG... from interpreted code
	// maybe a dot-import is easier to use?
	display, err := ir.Comp.ImportPackageOrError("display", "display")
	if err != nil {
		log.Print(err)
	}

	// Inject the stub "Display" function. declare a variable
	// instead of a function, because we want to later change
	// its value to the closure that holds a reference to msgReceipt
	ir.DeclVar("Display", nil, stubDisplay)

  adaptor := GoAdaptor{
		ir,
		display,
		nil,
	}
	adaptor.initRenderers()
  return &adaptor
}

// GetKernelInfo sends a kernel_info_reply message.
func (adaptor *GoAdaptor) GetKernelInfo() tk.KernelInfo {
  return tk.KernelInfo{
    ProtocolVersion:       tk.ProtocolVersion,
    Implementation:        "goIPyGophernotes",
    ImplementationVersion: Version,
    Banner:                fmt.Sprintf("Go kernel: goIPyGophernotes - v%s", Version),
    LanguageInfo:          tk.KernelLanguageInfo{
      Name:          "go",
      Version:       runtime.Version(),
      FileExtension: ".go",
    },
    HelpLinks: []tk.HelpLink{
      {Text: "Go", URL: "https://golang.org/"},
      {Text: "goIPyGophernotes", URL: "https://github.com/stephengaito/goIPythonKernelToolkit/kernels/goIPyGophernotes"},
    },
  }
}

func (adaptor *GoAdaptor) GetCodeWordCompletions(
  code string,
  cursorPos int,
) (int, int, []string) {
  
  // use the gomacro interpreter to find all matches to the word at the 
  // cursor. 
  //
  prefix, matches, _ := adaptor.ir.CompleteWords(code, cursorPos)
	partialWord        := gomacro.TailIdentifier(prefix)
  curStart           := len(prefix) - len(partialWord)
  curEnd             := cursorPos
  return curStart, curEnd, matches
}

func (adaptor *GoAdaptor) SetupDisplayCallback(receipt tk.MsgReceipt) {
  // inject the actual "Display" closure that displays multimedia data in Jupyter
	ir := adaptor.ir
	displayPlace := ir.ValueOf("Display")
	displayPlace.Set(reflect.ValueOf(receipt.PublishDisplayData))
}

func (adaptor *GoAdaptor) TeardownDisplayCallback() {
		// remove the closure before returning
	ir := adaptor.ir
	displayPlace := ir.ValueOf("Display")
  displayPlace.Set(reflect.ValueOf(stubDisplay))
}

// doEval evaluates the code in the interpreter. This function captures an 
// uncaught panic as well as the values of the last statement/expression. 
//
func (adaptor *GoAdaptor) EvaluateCode(
  execCounter int, 
  code        string,
) (rtnData tk.Data, err error) {
  ir := adaptor.ir
  
  // Capture a panic from the evaluation if one occurs and store it in the 
  // `err` return parameter. 
  //
  defer func() {
		if r := recover(); r != nil {
			var ok bool
			if err, ok = r.(error); !ok {
        rtnData = tk.Data{}
				err = errors.New(fmt.Sprint(r))
			}
		}
	}()

	// Prepare and perform the multiline evaluation.
	compiler := ir.Comp

	// Don't show the gomacro prompt.
	compiler.Options &^= base.OptShowPrompt

  // Don't swallow panics as they are recovered above and handled with a 
  // Jupyter `error` message instead. 
  //
  compiler.Options &^= base.OptTrapPanic

  // Reset the error line so that error messages correspond to the lines 
  // from the cell. 
  //
  compiler.Line = 0

  // Parse the input code (and don't perform gomacro's macroexpansion). 
  // These may panic but this will be recovered by the deferred recover() 
  // above so that the error may be returned instead. 
  //
  nodes := compiler.ParseBytes([]byte(code))
	srcAst := ast2.AnyToAst(nodes, "doEval")

  // If there is no srcAst then we must be evaluating nothing. The result 
  // must be nil then. 
  //
  if srcAst == nil {
		return tk.Data{}, nil
	}

  // Check if the last node is an expression. If the last node is not an 
  // expression then nothing is returned as a value. For example evaluating 
  // a function declaration shouldn't return a value but just have the side 
  // effect of declaring the function. 
  //
	// This is actually needed only for gomacro classic interpreter
	// (the fast interpreter already returns values only for expressions)
	// but retained for compatibility.
	var srcEndsWithExpr bool
	if len(nodes) > 0 {
		_, srcEndsWithExpr = nodes[len(nodes)-1].(ast.Expr)
	}

	// Compile the ast.
	compiledSrc := ir.CompileAst(srcAst)

	// Evaluate the code.
	results, types := ir.RunExpr(compiledSrc)
    
  // If the source ends with an expression, then the result of the execution 
  // is the value of the expression. In the event that all return values are 
  // nil, the result is also nil. 
  //
  if srcEndsWithExpr {

    // Count the number of non-nil values in the output. If they are all nil 
    // then the output is skipped. 
    //
    nonNilCount := 0
		values := make([]interface{}, len(results))
		for i, result := range results {
			val := basereflect.Interface(result)
			if val != nil {
				nonNilCount++
			}
			values[i] = val
		}

		if nonNilCount > 0 {
      data := adaptor.autoRenderResults(values, types)
      return data, nil
		}
	}

	return tk.Data{}, nil
}

// find and execute special commands in code, remove them from returned 
// string 
//
func (adaptor *GoAdaptor) EvaluateRemoveSpecialCommands(
  outerr tk.OutErr,
  code string,
) string {
  ir := adaptor.ir
  
	lines := strings.Split(code, "\n")
	stop := false
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) != 0 {
			switch line[0] {
			case '%':
				evalSpecialCommand(ir, outerr, line)
				lines[i] = ""
			case '$':
				evalShellCommand(ir, outerr, line)
				lines[i] = ""
			default:
				// if a line is NOT a special command,
				// stop processing special commands
				stop = true
			}
		}
		if stop {
			break
		}
	}
	return strings.Join(lines, "\n")
}

// execute special command. line must start with '%'
//
func evalSpecialCommand(ir *gomacro.Interp, outerr tk.OutErr, line string) {
	const help string = `
available special commands (%):
%help
%go111module {on|off}

execute shell commands ($): $command [args...]
example:
$ls -l
`

	args := strings.SplitN(line, " ", 2)
	cmd := args[0]
	arg := ""
	if len(args) > 1 {
		arg = args[1]
	}
	switch cmd {

	case "%go111module":
		if arg == "on" {
			ir.Comp.CompGlobals.Options |= base.OptModuleImport
		} else if arg == "off" {
			ir.Comp.CompGlobals.Options &^= base.OptModuleImport
		} else {
			panic(fmt.Errorf("special command %s: expecting a single argument 'on' or 'off', found: %q", cmd, arg))
		}
	case "%help":
		fmt.Fprint(outerr.Out, help)
	default:
		panic(fmt.Errorf("unknown special command: %q\n%s", line, help))
	}
}

// execute shell command. line must start with '$'
func evalShellCommand(ir *gomacro.Interp, outerr tk.OutErr, line string) {
	args := strings.Fields(line[1:])
	if len(args) <= 0 {
		return
	}

	var writersWG sync.WaitGroup
	writersWG.Add(2)

	cmd := exec.Command(args[0], args[1:]...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(fmt.Errorf("Command.StdoutPipe() failed: %v", err))
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		panic(fmt.Errorf("Command.StderrPipe() failed: %v", err))
	}

	go func() {
		defer writersWG.Done()
		io.Copy(outerr.Out, stdout)
	}()

	go func() {
		defer writersWG.Done()
		io.Copy(outerr.Err, stderr)
	}()

	err = cmd.Start()
	if err != nil {
		panic(fmt.Errorf("error starting command '%s': %v", line[1:], err))
	}

	err = cmd.Wait()
	if err != nil {
		panic(fmt.Errorf("error waiting for command '%s': %v", line[1:], err))
	}

	writersWG.Wait()
}
