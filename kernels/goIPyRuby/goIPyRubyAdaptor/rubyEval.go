package goIPyRubyAdaptor


// #cgo pkg-config: ruby
// #include <stdlib.h>
// #include <stdint.h>
// #include "rubyEval.h"
import "C"

import (
  "unsafe"
  //"fmt"
  
  tk "github.com/stephengaito/goIPythonKernelToolkit/goIPyKernel"
)


// Create a new Data object and store it in the IPyKernelStore.
//
// Return the GoUInt64 key to the new object in the IPyKernelStore.
//
//export GoIPyKernelData_New
func GoIPyKernelData_New() uint64 {
  return tk.TheObjectStore.Store(&tk.Data{
    Data:     make(tk.MIMEMap),
    Metadata: make(tk.MIMEMap),
  })
}

// Add the mimeType/dataValue pair to the Data map of the Data object.
//
// Takes the Data object at `objId` from the IPyKernelStore and adds the 
// mimeType/dataValue to the Data's Data map.
//
//export GoIPyKernelData_AddData
func GoIPyKernelData_AddData(
  objId         uint64,
  mimeTypePtr  *C.char,
  mimeTypeLen   C.int,
  dataValuePtr *C.char,
  dataValueLen  C.int,
) {
  anObj := tk.TheObjectStore.Get(objId)
  if anObj != nil {
    aDataObj := anObj.(*tk.Data)
    
    mimeType := C.GoStringN(mimeTypePtr, mimeTypeLen)
    if mimeType == tk.MIMETypePNG || mimeType == tk.MIMETypeJPEG {
      dataValue := C.GoBytes(unsafe.Pointer(dataValuePtr), dataValueLen)
      aDataObj.Data[mimeType] = dataValue
    } else {
      dataValue := C.GoStringN(dataValuePtr, dataValueLen)
      aDataObj.Data[mimeType] = dataValue
    }
  }
}

// Add the mimeType/metaKey/dataValue triple to the Metadata map of the Data object.
//
// Takes the Data object at `objId` from the IPyKernelStore and adds the 
// mimeType/metaKey/dataValue to the Data's Metadata map. 
//
//export GoIPyKernelData_AddMetadata
func GoIPyKernelData_AddMetadata(
  objId         uint64,
  mimeTypePtr  *C.char,
  mimeTypeLen   C.int,
  metaKeyPtr   *C.char,
  metaKeyLen    C.int,
  dataValuePtr *C.char,
  dataValueLen  C.int,
) {
  anObj := tk.TheObjectStore.Get(objId)
  if anObj != nil {
    aDataObj := anObj.(*tk.Data)
    
    mimeType := C.GoStringN(mimeTypePtr, mimeTypeLen)
    
    if aDataObj.Metadata[mimeType] == nil {
      aDataObj.Metadata[mimeType] = make(tk.MIMEMap)
    }
    aMimeMap  := aDataObj.Metadata[mimeType].(tk.MIMEMap)
    
    metaKey   := C.GoStringN(metaKeyPtr, metaKeyLen)
    dataValue := C.GoStringN(dataValuePtr, dataValueLen)
    
    aMimeMap[metaKey] = dataValue
  }
}

// A representation of the Ruby state.
//
// NOTE: Since Ruby is not reentrant, there can only be one Ruby instance.
//
type RubyState struct {
  // NO state to keep...
}

// Creates a running Ruby instance.
//
func CreateRubyState() *RubyState {
  C.startRuby()
  return &RubyState{}
}

// Stops the currently running Ruby instance.
//
func (rs *RubyState) DeleteRubyState() {
  C.stopRuby()
}

func (rs *RubyState) IsRubyRunning() bool {
  return C.isRubyRunning() != 0
}

// Return the Ruby version as a string
//
func (rs *RubyState) GetRubyVersion() string {
  return C.GoString(C.rubyVersion())
}

// Evaluate the String aGoStr in the (single) Ruby instance.
//
func (rs *RubyState) RubyEvaluateString(aGoStr string) tk.Data {
/*
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
*/
  return tk.Data{}
}
