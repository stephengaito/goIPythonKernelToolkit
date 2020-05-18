// +build cGoTests

package goIPyRubyAdaptor

import (
  "reflect"
  "testing"
  "github.com/stretchr/testify/assert"
  //"github.com/davecgh/go-spew/spew"
  tk "github.com/stephengaito/goIPythonKernelToolkit/goIPyKernel"
)

// assertions: https://godoc.org/github.com/stretchr/testify/assert
// prettyPrint: https://github.com/davecgh/go-spew

func TestGoIPyKernelData(t *testing.T) {
  objId := GoIPyKernelData_New()
  assert.Equal(t, tk.TheObjectStore.NextId, objId,
    "Object id should NextId",
  )
  
  GoAddMimeMapToDataObjTest(objId)
  
  anObj := tk.TheObjectStore.Get(objId)
  
  aDataObj := anObj.(*tk.Data)
  assert.NotNil(t, aDataObj.Data["MIMETest"], "Data object missing MIMETest")
  assert.Equalf(t, aDataObj.Data["MIMETest"], "some data",
    "Data object MIMETest has wrong value %s\n",
    aDataObj.Data["MIMETest"],
  )
  
  GoAddMimeMapToMetadataObjTest(objId)
  assert.NotNil(t, aDataObj.Metadata["MIMETest"],
    "Metadata object missing MIMETest",
  )
  aMimeMap := aDataObj.Metadata["MIMETest"].(tk.MIMEMap)
  assert.NotNil(t, aMimeMap["Width"],
    "Metadata object missing MIMETest.Width",
  )
  assert.Equalf(t, aMimeMap["Width"], "some data",
    "Metadata object MIMETest.Width has wrong value %s\n",
    aMimeMap["Width"],
  )
  
  GoAddJPEGMimeMapToDataObjTest(objId)
  assert.NotNil(t, aDataObj.Data["image/jpeg"],
    "Data object missing MIMEMap JPEG",
  )
  jpegObj    := aDataObj.Data["image/jpeg"];
  byteSlice := make([]byte, 0)
  assert.Equalf(t, reflect.TypeOf(jpegObj), reflect.TypeOf(byteSlice),
    "Data object MIMEMap JPEG has wrong type %T\n",
    jpegObj,
  )
  jpegSlice := jpegObj.([]byte)
  jpegLen   := len(jpegSlice)
  assert.Equalf(t, jpegLen, 10,
    "Data object MIMEMap JPEG has wrong length: %d\n",
    jpegLen,
  )
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
  assert.NotNil(t, aDataObj.Data["image/png"],
    "Data object missing MIMEMap PNG",
  )
  pngObj    := aDataObj.Data["image/png"];
  byteSlice  = make([]byte, 0)
  assert.Equalf(t, reflect.TypeOf(pngObj), reflect.TypeOf(byteSlice),
    "Data object MIMEMap PNG has wrong type %T\n",
    pngObj,
  )
  pngSlice := pngObj.([]byte)
  pngLen   := len(pngSlice)
  assert.Equalf(t, pngLen, 10,
    "Data object MIMEMap PNG has wrong length: %d\n",
    pngLen,
  )
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

func TestIPyKernelData_New(t *testing.T) {
  rubyState := CreateRubyState()
  
  codePath := "/lib/IPyRubyData.rb"
  IPyRubyDataCode, err := FSString(false, codePath)
  assert.NoErrorf(t, err,
    "Could not load file [%s]", codePath)
    
  _, err = rubyState.LoadRubyCode("IPyRubyData.rb", IPyRubyDataCode)
  assert.NoError(t, err, "Could not load IPyRubyData.rb")
 
  objId, err := 
    rubyState.LoadRubyCode("TestIPyKernelData_New", "IPyKernelData_New(nil)")
  
  assert.NoError(t, err, "Could not call TestIPyKernelData_New")
  anObj := tk.TheObjectStore.Get(objId)
  assert.NotNil(t, anObj, "Should have returned an empty Data object interface")
  aDataObj := anObj.(*tk.Data)
  assert.NotNil(t, aDataObj, "Should have returned an empty Data object")
  assert.Zero(t, len(aDataObj.Data),      "Should have an empty Data mimeMap")
  assert.Zero(t, len(aDataObj.Metadata),  "Should have an empty Metadata mimeMap")
  assert.Zero(t, len(aDataObj.Transient), "Should have an empty Transient mimeMap")
}

func TestIPyKernelData_AddData(t *testing.T) {
  rubyCode := `
    anObj = IPyKernelData_New(nil)
    IPyKernelData_AddData(anObj, MIMETypeText, "test text")
    anObj
  `
  rubyState := CreateRubyState()
  
  codePath := "/lib/IPyRubyData.rb"
  IPyRubyDataCode, err := FSString(false, codePath)
  assert.NoErrorf(t, err,
    "Could not load file [%s]", codePath)
    
  _, err = rubyState.LoadRubyCode("IPyRubyData.rb", IPyRubyDataCode)
  assert.NoError(t, err, "Could not load IPyRubyData.rb")
 
  objId, err := rubyState.LoadRubyCode("TestIPyKernelData_AddData", rubyCode)
  assert.NoError(t, err, "Could not call TestIPyKernelData_AddData")
  anObj := tk.TheObjectStore.Get(objId)
      
  assert.NotNil(t, anObj, "Should have returned an empty Data object interface")
  aDataObj := anObj.(*tk.Data)
  //spew.Dump(aDataObj)
  assert.NotNil(t, aDataObj, "Should have returned a Data object")
  assert.Equal(t, 1, len(aDataObj.Data),  "Should have a Data mimeMap with one entry")
  assert.Zero(t, len(aDataObj.Metadata),  "Should have an empty Metadata mimeMap")
  assert.Zero(t, len(aDataObj.Transient), "Should have an empty Transient mimeMap")
}

func TestIPyKernelData_AddMetadata(t *testing.T) {
  rubyCode := `
    anObj = IPyKernelData_New(nil)
    IPyKernelData_AddMetadata(anObj, MIMETypeText, "test", "test text")
    anObj
  `
  rubyState := CreateRubyState()
  
  codePath := "/lib/IPyRubyData.rb"
  IPyRubyDataCode, err := FSString(false, codePath)
  assert.NoErrorf(t, err,
    "Could not load file [%s]", codePath)
    
  _, err = rubyState.LoadRubyCode("IPyRubyData.rb", IPyRubyDataCode)
  assert.NoError(t, err, "Could not load IPyRubyData.rb")
 
  objId, err :=
    rubyState.LoadRubyCode("TestIPyKernelData_AddMetadata", rubyCode)
  assert.NoError(t, err, "Could not call TestIPyKernelData_AddMetadata")
  anObj := tk.TheObjectStore.Get(objId)
  
  assert.NotNil(t, anObj, "Should have returned a Data object interface")
  aDataObj := anObj.(*tk.Data)
  //spew.Dump(aDataObj)
  assert.NotNil(t, aDataObj, "Should have returned an empty Data object")
  assert.Zero(t, len(aDataObj.Data),  "Should have an empty Data mimeMap")
  assert.Equal(t, 1, len(aDataObj.Metadata),  "Should have a Metadata mimeMap with one entry")
  assert.Zero(t, len(aDataObj.Data),  "Should have an empty Metadata mimeMap")
  assert.Zero(t, len(aDataObj.Transient), "Should have an empty Transient mimeMap")
  
  assert.NotNil(t, aDataObj.Metadata["text/plain"],
    "aDataObj does not have MIMETypeText")
  aMimeMap := aDataObj.Metadata["text/plain"].(tk.MIMEMap)
  assert.NotNil(t, aMimeMap["test"],
    "aMimeMap should have test key")
  assert.Equal(t, aMimeMap["test"], "test text",
    "aMimeMap should have correct test value")
}

func TestMakeLastErrorData(t *testing.T) {
  // moved from lib/IPyRubyData_test.rb
  
  rubyTestCode := `savedErr = nil
    begin
      raise("This is silly")
    rescue
      savedErr = $!
    end
    MakeLastErrorData(savedErr, "This is a test")
`
  rubyState := CreateRubyState()
  
  codePath := "/lib/IPyRubyData.rb"
  IPyRubyDataCode, err := FSString(false, codePath)
  assert.NoErrorf(t, err,
    "Could not load file [%s]", codePath)
    
  _, err = rubyState.LoadRubyCode("IPyRubyData.rb", IPyRubyDataCode)
  assert.NoError(t, err, "Could not load IPyRubyData.rb")
 
  objId, err := rubyState.LoadRubyCode("TestMakeLastErrorData", rubyTestCode)
  assert.NoError(t, err, "Could not MakeLastErrorData")
  
  anObj := tk.TheObjectStore.Get(objId)
  
  assert.NotNil(t, anObj, "Should have returned a Data object interface")
  aDataObj := anObj.(*tk.Data)
  //spew.Dump(aDataObj)
  assert.NotNil(t, aDataObj, "Should have returned an empty Data object")
  assert.Equal(t, 4, len(aDataObj.Data),  "Should have an Metadata mimeMap with one entry")
  assert.Zero(t, len(aDataObj.Metadata),  "Should have an empty Metadata mimeMap")
  assert.Zero(t, len(aDataObj.Transient), "Should have an empty Transient mimeMap")

  lastErrData := aDataObj
  
  assert.NotNil(t, lastErrData.Data["ename"],
      "lastErrData does not have ename")
  assert.NotNil(t, lastErrData.Data["evalue"],
      "lastErrData does not have evalue")
  assert.NotNil(t, lastErrData.Data["traceback"],
      "lastErrData does not have traceback")
  assert.NotNil(t, lastErrData.Data["status"],
      "lastErrData does not have status")

  assert.Equal(t, lastErrData.Data["ename"], "ERROR",
      "lastErrData has incorrect evalue")
  assert.Equal(t, lastErrData.Data["evalue"], "This is silly",
      "lastErrData has incorrect evalue")
  assert.Equal(t, lastErrData.Data["traceback"].([]string)[0], "This is a test",
      "lastErrData has incorrect traceback")
  assert.Equal(t, lastErrData.Data["status"], "error",
      "lastErrData has incorrect status")
}