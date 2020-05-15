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

    // Test something
    //
    // Suite:   main
    // Fixture: main
    //
    func Test_IPyKernelDataCGoTest(t *testing.T) {      
      cGoTestMayBeError(t, "IPyKernelDataCGoTest", Go_IPyKernelDataCGoTest())
    }

    // no desc
    //
    // Suite:   main
    // Fixture: main
    //
    func Test_RubyStateCGoTest(t *testing.T) {      
      cGoTestMayBeError(t, "RubyStateCGoTest", Go_RubyStateCGoTest())
    }

  // end fixture: main

// end suite: main


