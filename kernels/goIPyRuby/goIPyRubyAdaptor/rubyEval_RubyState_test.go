package goIPyRubyAdaptor

import (
  "fmt"
  "testing"
   "github.com/stretchr/testify/assert"
   //"github.com/davecgh/go-spew/spew"
   //tk "github.com/stephengaito/goIPythonKernelToolkit/goIPyKernel"
)

// assertions: https://godoc.org/github.com/stretchr/testify/assert
// prettyPrint: https://github.com/davecgh/go-spew

func TestRubyState(t *testing.T) {

// Alas ruby has already been initialized in some other test.
//  fakeRubyState := &RubyState{}
//  if fakeRubyState.IsRubyRunning() {
//    t.Error("Ruby should NOT be running yet!")
//  }
 
  rubyState := CreateRubyState()
  assert.NotNil(t, rubyState,
    "rubyState should NOT be nil")

  assert.True(t, rubyState.IsRubyRunning(),
    "Ruby should be running now!" )

  rubyState = CreateRubyState()
  fmt.Printf("Ruby version: %s\n", rubyState.GetRubyVersion())
}

func TestLoadRubyCode(t *testing.T) {
  rubyState := CreateRubyState()
  
  _, err := rubyState.LoadRubyCode("goHelloWorldCode", "puts 'Hello world!'")
  assert.NoError(t,err, "Should have loaded hello world code")
    
  _, err = rubyState.LoadRubyCode("goBrokenCode", "this code is broken")
   assert.Error(t, err, "Should have not loaded broken code")
}

func TestLoadingIPyRubyDataCode(t *testing.T) {
  rubyState := CreateRubyState()
  
  codePath := "/lib/IPyRubyData.rb"
  IPyRubyDataCode, err := FSString(false, codePath)
  assert.NoErrorf(t, err,
    "Could not load file [%s]", codePath)

  _, err = rubyState.LoadRubyCode("IPyRubyData.rb", IPyRubyDataCode)
  assert.NoError(t, err, "Could not load IPyRubyData.rb")
  
  _, err =
    rubyState.LoadRubyCode("DoesIPyRubyEvalExist", "defined? IPyRubyEval")
  assert.NoError(t, err, "IPyRubyEval does not exist")  
}

func TestEvalRubyString(t *testing.T) {
  rubyState := CreateRubyState();
  
  dataObj := rubyState.GoEvalRubyString(
    "TestEvalRubyString1",
    "puts 'Hello TestEvalRubyString1'",
  )
  assert.NotNil(t, dataObj, "Should return a  non empty dataOjb")
  assert.NotNil(t, dataObj.Data["evalue"], "Should return an error")
  assert.Equal(t, dataObj.Data["evalue"], 
    "no return value from evalRubyString",
    "Should return correct error report",
  );
  assert.NotNil(t, dataObj.Data["state"], "Should be an error obj")
  assert.Equal(t, dataObj.Data["state"], "error", 
    "Should return correct error report",
  )
  assert.NotNil(t, dataObj.Data["ename"], "Should be an error obj")
  assert.Equal(t, dataObj.Data["ename"], "ERROR",
    "Should return correct error report",
  )
  assert.NotNil(t, dataObj.Data["traceback"], "Should be an error obj")
  assert.Equal(t, dataObj.Data["traceback"].([]string)[0], "GoEvalRubyString",
    "Should return correct error report",
  )
  //spew.Dump(dataObj)

  dataObj = rubyState.GoEvalRubyString(
    "TestEvalRubyString2",
    "a = 'Hello TestEvalRubyString2'",
  )
  assert.NotNil(t, dataObj, "Should return a  non empty dataOjb")
  assert.NotNil(t, dataObj.Data["evalue"], "Should return an error")
  assert.Equal(t, dataObj.Data["evalue"], 
    "no return value from evalRubyString",
    "Should return correct error report",
  );
  assert.NotNil(t, dataObj.Data["state"], "Should be an error obj")
  assert.Equal(t, dataObj.Data["state"], "error", 
    "Should return correct error report",
  )
  assert.NotNil(t, dataObj.Data["ename"], "Should be an error obj")
  assert.Equal(t, dataObj.Data["ename"], "ERROR",
    "Should return correct error report",
  )
  assert.NotNil(t, dataObj.Data["traceback"], "Should be an error obj")
  assert.Equal(t, dataObj.Data["traceback"].([]string)[0], "GoEvalRubyString",
    "Should return correct error report",
  )
  //spew.Dump(tk.TheObjectStore)
  //spew.Dump(dataObj)
}

func TestIPyRubyEval(t *testing.T) {
  rubyState := CreateRubyState();
  
  dataObj := rubyState.GoEvalRubyString(
    "TestIPyRubyEvalString1",
    "a = 'Hello TestIPyRubyEvalString1'",
  )
  assert.NotNil(t, dataObj, "Should return a non empty dataObj")
  
}

