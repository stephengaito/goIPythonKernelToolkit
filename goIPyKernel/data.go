package goIPyKernel

import (
  "fmt"
  "sync"
)

// MIMEMap holds data that can be presented in multiple formats. The keys are MIME types
// and the values are the data formatted with respect to its MIME type.
// All maps should contain at least a "text/plain" representation with a string value.
type MIMEMap map[string]interface{}


// DeepCopy makes a deep copy of the provided MIMEMap.
//
func (mm MIMEMap) DeepCopy() MIMEMap {
  newMM := make(MIMEMap, len(mm))
  for key, value := range mm {
    switch valueType := value.(type) {
    case string :
      newMM[key] = value
    case []byte :
      newMM[key] = value
    case []string :
      newValue := make([]string, len(value.([]string)))
      copy(newValue, value.([]string))
      newMM[key] = newValue
    case MIMEMap :
      newMM[key] = value.(MIMEMap).DeepCopy()
    default:
      panic(fmt.Sprintf("invalid MIMEMap type: %s ", valueType))
    }
  }
  return newMM
}

// Support an interface similar - but not identical - to the IPython 
// (canonical Jupyter kernel). See 
// http://ipython.readthedocs.io/en/stable/api/generated/IPython.display.html#IPython.display.display 
// for a good overview of the support types. 
//
const (
	MIMETypeHTML       = "text/html"
	MIMETypeJavaScript = "application/javascript"
	MIMETypeJPEG       = "image/jpeg"
	MIMETypeJSON       = "application/json"
	MIMETypeLatex      = "text/latex"
	MIMETypeMarkdown   = "text/markdown"
	MIMETypePNG        = "image/png"
	MIMETypePDF        = "application/pdf"
	MIMETypeSVG        = "image/svg+xml"
	MIMETypeText       = "text/plain"
)

// Data is the exact structure returned to Jupyter.
// It allows to fully specify how a value should be displayed.
type Data struct {
	Data      MIMEMap
	Metadata  MIMEMap
	Transient MIMEMap
}

// DeepCopy makes a deep copy of the provided Data structure.
//
func (data Data) DeepCopy() Data {
  return Data{
    Data:      data.Data.DeepCopy(),
    Metadata:  data.Metadata.DeepCopy(),
    Transient: data.Transient.DeepCopy(),
  }
}

// SyncedData is a synchronization wrapper around the Data structure. This 
// allows kernal adaptors to incrementally build up a Data structure one 
// element at a time rather than all at once. 
//
type SyncedData struct {
  Mutex  sync.RWMutex
  TKData Data
}

// Create a new Data object and store it in the IPyKernelStore.
//
// Return the GoUInt64 key to the new object in the IPyKernelStore.
//
func StoreData_New() uint64 {
  //fmt.Print("GoIPyKernelData_New\n")
  
  newObjId := TheObjectStore.Store(
    &SyncedData{
      TKData: Data{
        Data:      make(MIMEMap),
        Metadata:  make(MIMEMap),
        Transient: make(MIMEMap),
      },
    },
  )
  //fmt.Printf("  objId: %d\n", newObjId)
  return newObjId
}

// Delete an existing Data object from the IPyKernelStore.
//
func StoreData_Delete(objId uint64) {
  //fmt.Print("GoIPyKernelData_Delete\n")
  //fmt.Printf("  objId: %d\n", objId)
  
  TheObjectStore.Delete(objId)
}

// Add the mimeType/dataValue pair to the Data map of the Data object.
//
// Takes the Data object at `objId` from the IPyKernelStore and adds the 
// mimeType/dataValue to the Data's Data map.
//
func StoreData_AddStringData(
  objId     uint64,
  mimeType  string,
  dataValue string,
) {
  //fmt.Print("GoIPyKernelData_AddData\n")
  //fmt.Printf("  objId:     %d\n", objId)
  //fmt.Printf("  mimeType:  %s", mimeType)
  //fmt.Printf("  dataValue: %s\n", dataValue)
  
  anObj := TheObjectStore.Get(objId)
  
  if anObj != nil {
    aSyncedDataObj := anObj.(*SyncedData)
    aSyncedDataObj.Mutex.Lock()
    defer aSyncedDataObj.Mutex.Unlock()
    
    aSyncedDataObj.TKData.Data[mimeType] = dataValue
  }
}

// Add the mimeType/dataValue pair to the Data map of the Data object.
//
// Takes the Data object at `objId` from the IPyKernelStore and adds the 
// mimeType/dataValue to the Data's Data map.
//
func StoreData_AddBytesData(
  objId     uint64,
  mimeType  string,
  dataValue []byte,
) {
  //fmt.Print("GoIPyKernelData_AddData\n")
  //fmt.Printf("  objId:     %d\n", objId)
  //fmt.Printf("  mimeType:  %s", mimeType)
  //fmt.Printf("  dataValue: %s\n", dataValue)
  
  anObj := TheObjectStore.Get(objId)
  
  if anObj != nil {
    aSyncedDataObj := anObj.(*SyncedData)
    aSyncedDataObj.Mutex.Lock()
    defer aSyncedDataObj.Mutex.Unlock()
    
    aSyncedDataObj.TKData.Data[mimeType] = dataValue
  }
}

// Add the mimeType/dataValue pair to the Data map of the Data object.
//
// Takes the Data object at `objId` from the IPyKernelStore and adds one 
// traceback string to the Data's Data map. 
//
func StoreData_AppendTraceback(
  objId          uint64,
  tracebackValue string,
) {
  //fmt.Print("GoIPyKernelData_AppemdTraceback\n")
  //fmt.Printf("  objId:          %d\n", objId)
  //fmt.Printf("  tracebackValue: %s\n", tracebackValue)
  
  anObj := TheObjectStore.Get(objId)
  
  if anObj != nil {
    aSyncedDataObj := anObj.(*SyncedData)
    aSyncedDataObj.Mutex.Lock()
    defer aSyncedDataObj.Mutex.Unlock()
    
    if aSyncedDataObj.TKData.Data["traceback"] == nil {
      aSyncedDataObj.TKData.Data["traceback"] = make([]string, 0)
    }
    tracebackSlice := aSyncedDataObj.TKData.Data["traceback"].([]string)    
    //fmt.Printf("  tracebackSlice: %s\n", aDataObj.Data["traceback"])
    aSyncedDataObj.TKData.Data["traceback"] = 
      append(tracebackSlice, tracebackValue)
  }
}

// Add the mimeType/metaKey/dataValue triple to the Metadata map of the Data object.
//
// Takes the Data object at `objId` from the IPyKernelStore and adds the 
// mimeType/metaKey/dataValue to the Data's Metadata map. 
//
func StoreData_AddMetadata(
  objId     uint64,
  mimeType  string,
  metaKey   string,
  dataValue string,
) {
  //fmt.Print("GoIPyKernelData_AddMetadata\n")
  //fmt.Printf("  objId:     %d\n", objId)
  //fmt.Printf("  mimeType:  %s\n", mimeType)
  //fmt.Printf("  metaKey:   %s\n", metaKey)
  //fmt.Printf("  dataValue: %s\n", dataValue)
  
  anObj := TheObjectStore.Get(objId)
  if anObj != nil {
    aSyncedDataObj := anObj.(*SyncedData)
    aSyncedDataObj.Mutex.Lock()
    defer aSyncedDataObj.Mutex.Unlock()
    
    if aSyncedDataObj.TKData.Metadata[mimeType] == nil {
      aSyncedDataObj.TKData.Metadata[mimeType] = make(MIMEMap)
    }
    aMimeMap  := aSyncedDataObj.TKData.Metadata[mimeType].(MIMEMap)
    
    aMimeMap[metaKey] = dataValue
  }
}
