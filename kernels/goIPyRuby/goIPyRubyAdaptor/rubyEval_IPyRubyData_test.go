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

func TestGoIPyRubyData(t *testing.T) {
  objId := GoIPyRubyData_New()
  assert.Equal(t, tk.TheObjectStore.NextId, objId,
    "Object id should NextId",
  )
  
  GoAddMimeMapToDataObjTest(objId)
  
  anObj := tk.TheObjectStore.Get(objId)
  
  aSyncedDataObj := anObj.(*tk.SyncedData)
  assert.NotNil(t, aSyncedDataObj.TKData.Data["MIMETest"], "Data object missing MIMETest")
  assert.Equalf(t, aSyncedDataObj.TKData.Data["MIMETest"], "some data",
    "Data object MIMETest has wrong value %s\n",
    aSyncedDataObj.TKData.Data["MIMETest"],
  )
  
  GoAddMimeMapToMetadataObjTest(objId)
  assert.NotNil(t, aSyncedDataObj.TKData.Metadata["MIMETest"],
    "Metadata object missing MIMETest",
  )
  aMimeMap := aSyncedDataObj.TKData.Metadata["MIMETest"].(tk.MIMEMap)
  assert.NotNil(t, aMimeMap["Width"],
    "Metadata object missing MIMETest.Width",
  )
  assert.Equalf(t, aMimeMap["Width"], "some data",
    "Metadata object MIMETest.Width has wrong value %s\n",
    aMimeMap["Width"],
  )
  
  GoAddJPEGMimeMapToDataObjTest(objId)
  assert.NotNil(t, aSyncedDataObj.TKData.Data["image/jpeg"],
    "Data object missing MIMEMap JPEG",
  )
  jpegObj    := aSyncedDataObj.TKData.Data["image/jpeg"];
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
  assert.NotNil(t, aSyncedDataObj.TKData.Data["image/png"],
    "Data object missing MIMEMap PNG",
  )
  pngObj    := aSyncedDataObj.TKData.Data["image/png"];
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

func TestIPyRubyData_New(t *testing.T) {
  rubyState := CreateRubyState()
  
  codePath := "/lib/IPyRubyData.rb"
  IPyRubyDataCode, err := FSString(false, codePath)
  assert.NoErrorf(t, err,
    "Could not load file [%s]", codePath)
    
  _, err = rubyState.LoadRubyCode("IPyRubyData.rb", IPyRubyDataCode)
  assert.NoError(t, err, "Could not load IPyRubyData.rb")
 
  objId, err := 
    rubyState.LoadRubyCode("TestIPyRubyData_New", "IPyRubyData_New(nil)")
  
  assert.NoError(t, err, "Could not call TestIPyRubyData_New")
  anObj := tk.TheObjectStore.Get(uint64(objId))
  assert.NotNil(t, anObj, "Should have returned an empty Data object interface")
  aSyncedDataObj := anObj.(*tk.SyncedData)
  assert.NotNil(t, aSyncedDataObj, "Should have returned an empty Data object")
  assert.Zero(t, len(aSyncedDataObj.TKData.Data),      "Should have an empty Data mimeMap")
  assert.Zero(t, len(aSyncedDataObj.TKData.Metadata),  "Should have an empty Metadata mimeMap")
  assert.Zero(t, len(aSyncedDataObj.TKData.Transient), "Should have an empty Transient mimeMap")
}

func TestIPyRubyData_AddData(t *testing.T) {
  rubyCode := `
    anObj = IPyRubyData_New(nil)
    IPyRubyData_AddData(anObj, MIMETypeText, "test text")
    anObj
  `
  rubyState := CreateRubyState()
  
  codePath := "/lib/IPyRubyData.rb"
  IPyRubyDataCode, err := FSString(false, codePath)
  assert.NoErrorf(t, err,
    "Could not load file [%s]", codePath)
    
  _, err = rubyState.LoadRubyCode("IPyRubyData.rb", IPyRubyDataCode)
  assert.NoError(t, err, "Could not load IPyRubyData.rb")
 
  objId, err := rubyState.LoadRubyCode("TestIPyRubyData_AddData", rubyCode)
  assert.NoError(t, err, "Could not call TestIPyRubyData_AddData")
  anObj := tk.TheObjectStore.Get(uint64(objId))
      
  assert.NotNil(t, anObj, "Should have returned an empty Data object interface")
  aSyncedDataObj := anObj.(*tk.SyncedData)
  //spew.Dump(aSyncedDataObj)
  assert.NotNil(t, aSyncedDataObj, "Should have returned a Data object")
  assert.Equal(t, 1, len(aSyncedDataObj.TKData.Data),  "Should have a Data mimeMap with one entry")
  assert.Zero(t, len(aSyncedDataObj.TKData.Metadata),  "Should have an empty Metadata mimeMap")
  assert.Zero(t, len(aSyncedDataObj.TKData.Transient), "Should have an empty Transient mimeMap")
}

func TestIPyRubyData_AddMetadata(t *testing.T) {
  rubyCode := `
    anObj = IPyRubyData_New(nil)
    IPyRubyData_AddMetadata(anObj, MIMETypeText, "test", "test text")
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
    rubyState.LoadRubyCode("TestIPyRubyData_AddMetadata", rubyCode)
  assert.NoError(t, err, "Could not call TestIPyRubyData_AddMetadata")
  anObj := tk.TheObjectStore.Get(uint64(objId))
  
  assert.NotNil(t, anObj, "Should have returned a Data object interface")
  aSyncedDataObj := anObj.(*tk.SyncedData)
  //spew.Dump(aSyncedDataObj)
  assert.NotNil(t, aSyncedDataObj, "Should have returned an empty Data object")
  assert.Zero(t, len(aSyncedDataObj.TKData.Data),  "Should have an empty Data mimeMap")
  assert.Equal(t, 1, len(aSyncedDataObj.TKData.Metadata),  "Should have a Metadata mimeMap with one entry")
  assert.Zero(t, len(aSyncedDataObj.TKData.Data),  "Should have an empty Metadata mimeMap")
  assert.Zero(t, len(aSyncedDataObj.TKData.Transient), "Should have an empty Transient mimeMap")
  
  assert.NotNil(t, aSyncedDataObj.TKData.Metadata["text/plain"],
    "aDataObj does not have MIMETypeText")
  aMimeMap := aSyncedDataObj.TKData.Metadata["text/plain"].(tk.MIMEMap)
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
  
  anObj := tk.TheObjectStore.Get(uint64(objId))
  
  assert.NotNil(t, anObj, "Should have returned a Data object interface")
  aSyncedDataObj := anObj.(*tk.SyncedData)
  //spew.Dump(aSyncedDataObj)
  assert.NotNil(t, aSyncedDataObj, "Should have returned an empty Data object")
  assert.Equal(t, 4, len(aSyncedDataObj.TKData.Data),  "Should have an Metadata mimeMap with one entry")
  assert.Zero(t, len(aSyncedDataObj.TKData.Metadata),  "Should have an empty Metadata mimeMap")
  assert.Zero(t, len(aSyncedDataObj.TKData.Transient), "Should have an empty Transient mimeMap")

  lastErrData := aSyncedDataObj
  
  assert.NotNil(t, lastErrData.TKData.Data["ename"],
      "lastErrData does not have ename")
  assert.NotNil(t, lastErrData.TKData.Data["evalue"],
      "lastErrData does not have evalue")
  assert.NotNil(t, lastErrData.TKData.Data["traceback"],
      "lastErrData does not have traceback")
  assert.NotNil(t, lastErrData.TKData.Data["status"],
      "lastErrData does not have status")

  assert.Equal(t, lastErrData.TKData.Data["ename"], "ERROR",
      "lastErrData has incorrect evalue")
  assert.Equal(t, lastErrData.TKData.Data["evalue"], "This is silly",
      "lastErrData has incorrect evalue")
  assert.Equal(t, lastErrData.TKData.Data["traceback"].([]string)[0], "This is a test",
      "lastErrData has incorrect traceback")
  assert.Equal(t, lastErrData.TKData.Data["status"], "error",
      "lastErrData has incorrect status")
}