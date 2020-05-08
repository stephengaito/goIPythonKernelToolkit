package goIPyRubyAdaptor

// #include <stdlib.h>
// #include "rubyEval.h"
import "C"

import (
  "unsafe"
)

type GOIPythonReturn struct {
  MimeType string,
  Value    string,
}

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
