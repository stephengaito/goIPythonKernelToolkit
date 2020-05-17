// +buid cGoTests

// Some tests of the Ruby interface using cGoTests

#include <ruby.h>
#include "_cgo_export.h"
#include "goIPyRubyAdaptorCGoTests.h"
#include "cGoTests.h"


/// \brief Should only load the helloWorldCode once
///
char *LoadHelloWorldCodeCGoTest(void *data) {
  startRuby();

  char *errMesg = loadRubyCode("helloWorldCode", "puts 'Hello World!'");
  cGoTest_Nil("Should have no error message", errMesg);
  
  cGoTest(
    "Should load the helloWorldCode",
    isRubyCodeLoaded("helloWorldCode")
  );

  errMesg = loadRubyCode("helloWorldCode", "puts 'Hello World!'");
  cGoTest_Nil("Should have no error message", errMesg);
  
  cGoTest(
    "Should load the helloWorldCode",
    isRubyCodeLoaded("helloWorldCode")
  );

  return NULL;
}

/// \brief Should fail to load the brokenCode
///
char *LoadBrokenCodeCGoTest(void *data) {
  startRuby();

  char *errMesg = loadRubyCode("borkenCode", "this code is broken");
  cGoTest_NotNil("Should have error message", errMesg);
  printf("error message: [%s]\n", errMesg);
  cGoTest_StrContains("Should contain the word broken", errMesg, "broken")
  cGoTest(
    "Should not load the brokenCode",
    ! isRubyCodeLoaded("brokenCode")
  );

  return NULL;
}