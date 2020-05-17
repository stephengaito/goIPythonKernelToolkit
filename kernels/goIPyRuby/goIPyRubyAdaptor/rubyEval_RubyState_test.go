package goIPyRubyAdaptor

import (
  "fmt"
  "testing"
   "github.com/stretchr/testify/assert"
)

// assertions: https://godoc.org/github.com/stretchr/testify/assert

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
  
  assert.NoError(t,
    rubyState.LoadRubyCode("goHelloWorldCode", "puts 'Hello world!'"),
    "Should have loaded hello world code")
    
   assert.Error(t,
    rubyState.LoadRubyCode("goBrokenCode", "this code is broken"),
    "Should have not loaded broken code")
}

func TestLoadingIPyRubyDataCode(t *testing.T) {
  rubyState := CreateRubyState()
  
  codePath := "/lib/IPyRubyData.rb"
  IPyRubyDataCode, err := FSString(false, codePath)
  assert.NoErrorf(t, err,
    "Could not load file [%s]", codePath)
    
  assert.NoError(t, 
    rubyState.LoadRubyCode("IPyRubyData.rb", IPyRubyDataCode),
    "Could not load IPyRubyData.rb")
  
  assert.NoError(t,
    rubyState.LoadRubyCode("DoesIPyRubyEvalExist", "defined? IPyRubyEval"),
    "IPyRubyEval does not exist")  
}