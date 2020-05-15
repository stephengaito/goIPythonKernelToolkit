// +build cGoTests

// GoLang level wrappers of the ANSI-C tests in the goIPyRubyAdaptor
// GoLang Package. 
//
// This *should* be located in goIPyRubyAdaptorCGoTests_test.go...
// ... unfortunately `go test` forbids the use of cgo...
// ... so we need to maintain this addition level of code indirection.
//
// Package description:
//   goIPyRubyAdaptor ANSI-C tests
//
// This file is automatically (re)generated changes made to this file will 
// be lost. 
//
package goIPyRubyAdaptor

// #include "goIPyRubyAdaptorCGoTests.h"
import "C"

import (
)


// begin suite: main

  // begin fixture: main

    // Test something
    //
    // Suite:   main
    // Fixture: main
    //
    func Go_IPyKernelDataCGoTest() error {
      cGoTestSuite("main", "Main (default) TestSuite")
      cGoTestFixture("main", "Main (default) Fixture in Main Suite")
      
      cGoTestStart("IPyKernelDataCGoTest", "Test something")
      defer cGoTestFinish("IPyKernelDataCGoTest")


      data := C.nullSetup()


      
      return cGoTestPossibleError(C.IPyKernelDataCGoTest(data))
    }

    // no desc
    //
    // Suite:   main
    // Fixture: main
    //
    func Go_RubyStateCGoTest() error {
      cGoTestSuite("main", "Main (default) TestSuite")
      cGoTestFixture("main", "Main (default) Fixture in Main Suite")
      
      cGoTestStart("RubyStateCGoTest", "no desc")
      defer cGoTestFinish("RubyStateCGoTest")


      data := C.nullSetup()


      
      return cGoTestPossibleError(C.RubyStateCGoTest(data))
    }

  // end fixture: main

// end suite: main


