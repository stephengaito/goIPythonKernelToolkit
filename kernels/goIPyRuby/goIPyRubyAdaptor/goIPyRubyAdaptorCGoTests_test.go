// +build cGoTests

// GoLang level tests for the goIPyRubyAdaptor ANSI-C code
//
// Package description:
//   goIPyRubyAdaptor ANSI-C tests
//
// This file is automatically (re)generated changes made to this file will 
// be lost. 

package goIPyRubyAdaptor

import (
  "testing"
)


// begin suite: main

  // begin fixture: main

    // Should evaluate some ruby code
    //
    // Suite:   main
    // Fixture: main
    //
    func Test_EvalRubyStringCGoTest(t *testing.T) {      
      cGoTestMayBeError(t, "EvalRubyStringCGoTest", Go_EvalRubyStringCGoTest())
    }

    // Should fail to load the brokenCode
    //
    // Suite:   main
    // Fixture: main
    //
    func Test_LoadBrokenCodeCGoTest(t *testing.T) {      
      cGoTestMayBeError(t, "LoadBrokenCodeCGoTest", Go_LoadBrokenCodeCGoTest())
    }

    // Should only load the helloWorldCode once
    //
    // Suite:   main
    // Fixture: main
    //
    func Test_LoadHelloWorldCodeCGoTest(t *testing.T) {      
      cGoTestMayBeError(t, "LoadHelloWorldCodeCGoTest", Go_LoadHelloWorldCodeCGoTest())
    }

  // end fixture: main

// end suite: main


