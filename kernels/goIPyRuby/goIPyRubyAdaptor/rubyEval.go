package goIPyRubyAdaptor


// #cgo pkg-config: ruby
// #include <stdlib.h>
// #include <stdint.h>
// #include "rubyEval.h"
import "C"

import (
  "unsafe"
  "errors"
  //"fmt"
  
  tk "github.com/stephengaito/goIPythonKernelToolkit/goIPyKernel"
)


// Create a new Data object and store it in the IPyKernelStore.
//
// Return the GoUInt64 key to the new object in the IPyKernelStore.
//
//export GoIPyKernelData_New
func GoIPyKernelData_New() uint64 {
  //fmt.Print("GoIPyKernelData_New\n")
  newObjId := tk.TheObjectStore.Store(&tk.Data{
    Data:      make(tk.MIMEMap),
    Metadata:  make(tk.MIMEMap),
    Transient: make(tk.MIMEMap),
  })
  //fmt.Printf("  objId:       %d\n", newObjId)
  return newObjId
}

// Delete an existing Data object from the IPyKernelStore.
//
//export GoIPyKernelData_Delete
func GoIPyKernelData_Delete(objId uint64) {
  //fmt.Print("GoIPyKernelData_Delete\n")
  //fmt.Printf("  objId:       %d\n", objId)
  tk.TheObjectStore.Delete(objId)
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
  //fmt.Print("GoIPyKernelData_AddData\n")
  //fmt.Printf("  objId:       %d\n", objId)
  anObj := tk.TheObjectStore.Get(objId)
  
  if anObj != nil {
    aDataObj := anObj.(*tk.Data)
    
    mimeType := C.GoStringN(mimeTypePtr, mimeTypeLen)
    //fmt.Printf("  mimeType:  %s", mimeType)
    if mimeType == tk.MIMETypePNG || mimeType == tk.MIMETypeJPEG {
      dataValue := C.GoBytes(unsafe.Pointer(dataValuePtr), dataValueLen)
      //fmt.Printf("  dataValue: %s\n", dataValue)
      aDataObj.Data[mimeType] = dataValue
    } else {
      dataValue := C.GoStringN(dataValuePtr, dataValueLen)
      //fmt.Printf("  dataValue: %s\n", dataValue)
      aDataObj.Data[mimeType] = dataValue
    }
  }
}

// Add the mimeType/dataValue pair to the Data map of the Data object.
//
// Takes the Data object at `objId` from the IPyKernelStore and adds one 
// traceback string to the Data's Data map. 
//
//export GoIPyKernelData_AppendTraceback
func GoIPyKernelData_AppendTraceback(
  objId              uint64,
  tracebackValuePtr *C.char,
  tracebackValueLen  C.int,
) {
  //fmt.Print("GoIPyKernelData_AppemdTraceback\n")
  //fmt.Printf("  objId:       %d\n", objId)
  anObj := tk.TheObjectStore.Get(objId)
  
  if anObj != nil {
    aDataObj := anObj.(*tk.Data)
    
    tracebackValue := C.GoStringN(tracebackValuePtr, tracebackValueLen)
    if aDataObj.Data["traceback"] == nil {
      aDataObj.Data["traceback"] = make([]string, 0)
    }
    tracebackSlice := aDataObj.Data["traceback"].([]string)    
    aDataObj.Data["traceback"] = 
      append(tracebackSlice, tracebackValue)
    
    //fmt.Printf("  tracebackValue: %s\n", tracebackValue)
    //fmt.Printf("  tracebackSlice: %s\n", aDataObj.Data["traceback"])
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
  //fmt.Print("GoIPyKernelData_AddMetadata\n")
  //fmt.Printf("  objId:       %d\n", objId)
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

    //fmt.Printf("  mimeType:  %s\n", mimeType)
    //fmt.Printf("  metaKey:   %s\n", metaKey)
    //fmt.Printf("  dataValue: %s\n", dataValue)
    
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

// Returns true if the Ruby virtual machine is running.
//
func (rs *RubyState) IsRubyRunning() bool {
  return C.isRubyRunning() != 0
}

// Load the Ruby code named `rubyCodeName` from the contents of the string 
// `rubyCode`.
//
// Returns an int64 and (possibly) a Go error structure.
//
// If the `rubyCode` returns any Integers, that integer is returned in the 
// int64. Otherwise the int64 is zero. 
//
// If any errors occur, a string description is wrapped in a Go error 
// structure. 
//
func (rs *RubyState) LoadRubyCode(
  rubyCodeName string,
  rubyCode     string,
) (int64, error) {
  rubyCodeCStr     := C.CString(rubyCode)
  defer C.free(unsafe.Pointer(rubyCodeCStr))
  rubyCodeNameCStr := C.CString(rubyCodeName)
  defer C.free(unsafe.Pointer(rubyCodeNameCStr))
  
  result := C.loadRubyCode(rubyCodeNameCStr, rubyCodeCStr)
  if result == nil {
    return 0, errors.New("No LoadRubyCode result structure returned")
  }
  defer C.FreeLoadRubyCodeReturn(result)
  
  if result.errMesg != nil {
    errMesg := C.GoString(result.errMesg)
    return int64(result.objId), errors.New(errMesg)
  }
  return int64(result.objId), nil
}

// Returns true if the Ruby code named `rubyCodeName` has already been 
// loaded. 
//
func (rs *RubyState) IsRubyCodeLoaded(rubyCodeName string) bool {
  rubyCodeNameCStr := C.CString(rubyCodeName)
  defer C.free(unsafe.Pointer(rubyCodeNameCStr))

  return C.isRubyCodeLoaded(rubyCodeNameCStr) != 0
}

// Return the Ruby version as a string
//
func (rs *RubyState) GetRubyVersion() string {
  return C.GoString(C.rubyVersion())
}

// Evaluate the String aGoStr in the (single) Ruby instance.
//
func (rs *RubyState) GoEvalRubyString(
  rubyCodeName, rubyCodeStr string,
) *tk.Data {  
  rubyCodeNameCStr := C.CString(rubyCodeName)
  defer C.free(unsafe.Pointer(rubyCodeNameCStr))
  
  rubyCodeCStr := C.CString(rubyCodeStr)
  defer C.free(unsafe.Pointer(rubyCodeCStr))

  objId := uint64(C.evalRubyString(rubyCodeNameCStr, rubyCodeCStr))
  if objId == 0 {
    return &tk.Data{
      Data: tk.MIMEMap{
        "ename":     "ERROR",
        "evalue":    "no return value from evalRubyString",
        "traceback": []string{ "GoEvalRubyString" },
        "state":     "error",
      },
      Metadata: tk.MIMEMap{},
      Transient: tk.MIMEMap{},
    }
  }

  storedObj := tk.TheObjectStore.Get(objId)
  if storedObj == nil {
    return &tk.Data{
      Data: tk.MIMEMap{
        "ename":     "ERROR",
        "evalue":    "no data object in the object store",
        "traceback": []string { "GoEvalRubyString" },
        "state":     "error",
      },
      Metadata: tk.MIMEMap{},
      Transient: tk.MIMEMap{},
    }  
  }
  return storedObj.(*tk.Data)
}
