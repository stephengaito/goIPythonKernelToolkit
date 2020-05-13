package goIPyRubyAdaptor


// #cgo pkg-config: ruby
// #include <stdlib.h>
// #include "rubyEval.h"
import "C"


import (
  //"unsafe"
  "fmt"
  
  tk "github.com/stephengaito/goIPythonKernelToolkit/goIPyKernel"
)

const (
	// Version defines the goIPyGophernotes version.
	Version string = "1.0.0"
)

type GOIPythonReturn struct {
  MimeType string
  Value    string
}

/*
func (r *RubyState) evalString(aGoStr string) {
  const char* aCStr = C.CString(aGoStr)
  defer C.free(unsafe.Pointer(aCStr))
  
  returnValue, err := C.evalString(aCStr)
  if err != nil {
    // do something
  }
  defer C.freeIPythonReturn(returnValue)
  
  return &GOIPythonReturn{
    MimeType: C.GOString(returnValue.mimeType),
    Value:    C.GOString(returnValue.value),
  }
}
*/
type GoAdaptor struct {

}

func NewGoAdaptor() *GoAdaptor {
  return &GoAdaptor{}
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
      Version:       C.GoString(C.rubyVersion()),
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
