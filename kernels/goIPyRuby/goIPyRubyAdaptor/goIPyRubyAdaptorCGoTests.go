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

    // GoIPyKernelData_New should return a new object id.
    //
    // Suite:   main
    // Fixture: main
    //
    func Go_IPyKernelDataCGoTest() error {
      cGoTestSuite("main", "Main (default) TestSuite")
      cGoTestFixture("main", "Main (default) Fixture in Main Suite")
      
      cGoTestStart("IPyKernelDataCGoTest", "GoIPyKernelData_New should return a new object id.")
      defer cGoTestFinish("IPyKernelDataCGoTest")


      data := C.nullSetup()


      
      return cGoTestPossibleError(C.IPyKernelDataCGoTest(data))
    }

    // Should fail to load the brokenCode
    //
    // Suite:   main
    // Fixture: main
    //
    func Go_LoadBrokenCodeCGoTest() error {
      cGoTestSuite("main", "Main (default) TestSuite")
      cGoTestFixture("main", "Main (default) Fixture in Main Suite")
      
      cGoTestStart("LoadBrokenCodeCGoTest", "Should fail to load the brokenCode")
      defer cGoTestFinish("LoadBrokenCodeCGoTest")


      data := C.nullSetup()


      
      return cGoTestPossibleError(C.LoadBrokenCodeCGoTest(data))
    }

    // Should only load the helloWorldCode once
    //
    // Suite:   main
    // Fixture: main
    //
    func Go_LoadHelloWorldCodeCGoTest() error {
      cGoTestSuite("main", "Main (default) TestSuite")
      cGoTestFixture("main", "Main (default) Fixture in Main Suite")
      
      cGoTestStart("LoadHelloWorldCodeCGoTest", "Should only load the helloWorldCode once")
      defer cGoTestFinish("LoadHelloWorldCodeCGoTest")


      data := C.nullSetup()


      
      return cGoTestPossibleError(C.LoadHelloWorldCodeCGoTest(data))
    }

  // end fixture: main

// end suite: main


