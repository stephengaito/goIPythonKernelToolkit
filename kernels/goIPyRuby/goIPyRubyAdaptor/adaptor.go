//go:generate esc -o rubyCode.go -pkg goIPyRubyAdaptor lib/IPyRubyData.rb
//go:generate cGoTestGenerator goIPyRubyAdaptor goIPyRubyAdaptor ANSI-C tests

package goIPyRubyAdaptor

import (
  //"unsafe"
  "fmt"
  
  tk "github.com/stephengaito/goIPythonKernelToolkit/goIPyKernel"
)

const (
	// Version defines the goIPyGophernotes version.
	Version string = "1.0.0"
)

// GoAdaptor represents any state required by the adaptor.
///
type GoAdaptor struct {
  // The Ruby State
  Ruby *RubyState
}

// Create a new adaptor.
//
func NewGoAdaptor() *GoAdaptor {
  return &GoAdaptor{
    Ruby: CreateRubyState(),
  }
}

// GetKernelInfo returns the KernelInfo for this kernel implementation.
//
func (adaptor *GoAdaptor) GetKernelInfo() tk.KernelInfo {
  return tk.KernelInfo{
    ProtocolVersion:       tk.ProtocolVersion,
    Implementation:        "goIPyRuby",
    ImplementationVersion: Version,
    Banner:                fmt.Sprintf("Go kernel: goIPyRuby - v%s", Version),
    LanguageInfo:          tk.KernelLanguageInfo{
      Name:          "ruby",
      Version:       adaptor.Ruby.GetRubyVersion(),
      FileExtension: ".rb",
    },
    HelpLinks: []tk.HelpLink{
      {Text: "Ruby", URL: "https://golang.org/"},
      {Text: "goIPyRuby", URL: "https://github.com/stephengaito/goIPythonKernelToolkit/kernels/goIPyRuby"},
    },
  }
}
  
// Get the possible completions for the word at cursorPos in the code. 
//
func (adaptor *GoAdaptor) GetCodeWordCompletions(
  code string,
  cursorPos int,
) (int, int, []string) {
  return 0, 0, make([]string, 0)
}

// Setup the Display callback by recording the msgReceipt information
// for later use by what ever callback implements the "Display" function. 
//
func (adaptor *GoAdaptor) SetupDisplayCallback(receipt tk.MsgReceipt) {
}
  
// Teardown the Display callback by removing the current msgReceipt
// information and setting things back to what ever default the 
// implementation uses.
//
func (adaptor *GoAdaptor) TeardownDisplayCallback() {
}
  
// Evaluate (and remove) any implmenation specific special commands BEFORE 
// the code gets evaluated by the interpreter. The `outErr` variable 
// contains the stdOut and stdErr which can be used to capture the stdOut 
// and stdErr of any external commands run by these special commands. 
//
func (adaptor *GoAdaptor) EvaluateRemoveSpecialCommands(
  outErr tk.OutErr,
  code string,
) string {
  return code
}

// Evaluate the code and return the results as a Data object.
//
func (adaptor *GoAdaptor) EvaluateCode(code string) (rtnData tk.Data, err error) {
  return tk.Data{}, nil
}
