// +build cGoTests

package goIPyRubyAdaptor

import (
  "reflect"
  "testing"
  
  tk "github.com/stephengaito/goIPythonKernelToolkit/goIPyKernel"
)

func TestGoIPyKernelData(t *testing.T) {
  objId := GoIPyKernelData_New()
  if objId != tk.TheObjectStore.NextId {
    t.Errorf(
      "Object id should be %d but is: %d\n",
      tk.TheObjectStore.NextId,
      objId,
    )
  }
  
  GoAddMimeMapToDataObjTest(objId)
  anObj := tk.TheObjectStore.Get(objId)
  aDataObj := anObj.(*tk.Data)
  if aDataObj.Data["MIMETest"] == nil {
    t.Error("Data object missing MIMETest")
  }
  if aDataObj.Data["MIMETest"] != "some data" {
    t.Errorf(
      "Data object MIMETest has wrong value %s\n",
      aDataObj.Data["MIMETest"],
    )
  }
  
  GoAddMimeMapToMetadataObjTest(objId)
  if aDataObj.Metadata["MIMETest"] == nil {
    t.Error("Metadata object missing MIMETest")
  }
  aMimeMap := aDataObj.Metadata["MIMETest"].(tk.MIMEMap)
  if aMimeMap["Width"] == nil {
    t.Error("Metadata object missing MIMETest.Width")
  }
  if aMimeMap["Width"] != "some data" {
    t.Errorf(
      "Metadata object MIMETest.Width has wrong value %s\n",
      aMimeMap["Width"],
    )
  }
  
  GoAddJPEGMimeMapToDataObjTest(objId)
  if aDataObj.Data["image/jpeg"] == nil {
    t.Error("Data object missing MIMEMap JPEG")
  }
  jpegObj    := aDataObj.Data["image/jpeg"];
  byteSlice := make([]byte, 0)
  if reflect.TypeOf(jpegObj) != reflect.TypeOf(byteSlice) {
    t.Errorf(
      "Data object MIMEMap JPEG has wrong type %T\n",
      jpegObj,
    )
  }
  jpegSlice := jpegObj.([]byte)
  jpegLen   := len(jpegSlice)
  if jpegLen != 10 {
    t.Errorf(
      "Data object MIMEMap JPEG has wrong length: %d\n",
      jpegLen,
    )
  }
  if jpegSlice[0] != 's' ||
     jpegSlice[1] != 'o' ||
     jpegSlice[2] != 'm' ||
     jpegSlice[3] != 'e' ||
     jpegSlice[4] != 0   ||
     jpegSlice[5] != 'd' ||
     jpegSlice[6] != 'a' ||
     jpegSlice[7] != 't' ||
     jpegSlice[8] != 'a' ||
     jpegSlice[9] != 0 {
    t.Error("Data object MIMEMap PJPEG has wrong content")
  }
  
  GoAddPNGMimeMapToDataObjTest(objId)
  if aDataObj.Data["image/png"] == nil {
    t.Error("Data object missing MIMEMap PNG")
  }
  pngObj    := aDataObj.Data["image/png"];
  byteSlice  = make([]byte, 0)
  if reflect.TypeOf(pngObj) != reflect.TypeOf(byteSlice) {
    t.Errorf(
      "Data object MIMEMap PNG has wrong type %T\n",
      pngObj,
    )
  }
  pngSlice := pngObj.([]byte)
  pngLen   := len(pngSlice)
  if pngLen != 10 {
    t.Errorf(
      "Data object MIMEMap PNG has wrong length: %d\n",
      pngLen,
    )
  }
  if pngSlice[0] != 's' ||
     pngSlice[1] != 'o' ||
     pngSlice[2] != 'm' ||
     pngSlice[3] != 'e' ||
     pngSlice[4] != 0   ||
     pngSlice[5] != 'd' ||
     pngSlice[6] != 'a' ||
     pngSlice[7] != 't' ||
     pngSlice[8] != 'a' ||
     pngSlice[9] != 0 {
    t.Error("Data object MIMEMap PNG has wrong content")
  }
}