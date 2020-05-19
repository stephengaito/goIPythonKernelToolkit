// +buid cGoTests

/// \file
/// \brief Some tests of the Ruby interface using cGoTests

#include <ruby.h>
#include "_cgo_export.h"
#include "goIPyRubyAdaptorCGoTests.h"
#include "cGoTests.h"


/// \brief Should only load the helloWorldCode once
///
char *LoadHelloWorldCodeCGoTest(void *data) {
  startRuby();

  LoadRubyCodeReturn *result = loadRubyCode("helloWorldCode",
    "puts 'Hello LoadHelloWorldCodeCGoTest!'");
  cGoTest_NotNil_MayFail("Should have returned a result", result);
  cGoTest_Nil("Should have no error message", result->errMesg);
  result = FreeLoadRubyCodeReturn(result);
  
  cGoTest(
    "Should load the helloWorldCode",
    isRubyCodeLoaded("helloWorldCode")
  );

  result = loadRubyCode("helloWorldCode", "puts 'Hello World!'");
  cGoTest_NotNil_MayFail("Should have returned a result", result);
  cGoTest_Nil("Should have no error message", result->errMesg);
  result = FreeLoadRubyCodeReturn(result);
  
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

  LoadRubyCodeReturn *result =
    loadRubyCode("brokenCode", "this code is broken");
  cGoTest_NotNil_MayFail("Should have returned a result", result);
  cGoTest_NotNil("Should have error message", result->errMesg);
  printf("error message: [%s]\n", result->errMesg);
  cGoTest_StrContains(
    "Should contain the word broken",
    result->errMesg, "broken"
  )
  cGoTest(
    "Should not load the brokenCode",
    ! isRubyCodeLoaded("brokenCode")
  );

  return NULL;
}

/// \brief Should evaluate some ruby code
///
char *EvalRubyStringCGoTest(void *data) {
  startRuby();
  
  uint64_t result = evalRubyString(
    "evalRubyStringCGoTest",
    "puts 'Hello EvalRubyStringCGoTest!'"
  );
  cGoTest_UIntNotEquals(
    "Should not be the zero object in the object store",
    result, 0
  );
  return NULL;
}

