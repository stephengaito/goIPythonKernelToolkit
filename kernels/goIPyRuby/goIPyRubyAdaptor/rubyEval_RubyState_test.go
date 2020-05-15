package goIPyRubyAdaptor

import (
  "fmt"
  "testing"
)

func TestRubyState(t *testing.T) {

// Alas ruby has already been initialized in some other test.
//  fakeRubyState := &RubyState{}
//  if fakeRubyState.IsRubyRunning() {
//    t.Error("Ruby should NOT be running yet!")
//  }
 
  rubyState := CreateRubyState()
  if rubyState == nil {
    t.Error("rubyState should NOT be nil")
  }
  if ! rubyState.IsRubyRunning() {
    t.Error("Ruby should be running now!")
  }

  rubyState = CreateRubyState()
  fmt.Printf("Ruby version: %s\n", rubyState.GetRubyVersion())
}