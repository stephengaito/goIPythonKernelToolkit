// +buid cGoTests

// Some tests of the Ruby interface using cGoTests

#include <ruby.h>
#include "_cgo_export.h"
#include "goIPyRubyAdaptorCGoTests.h"
#include "cGoTests.h"

char *RubyStateCGoTest(void *data) {
  startRuby();
  rb_set_errinfo(Qnil);
  int rubyFailed;
  VALUE result;
  result = rb_eval_string_protect("require 'pp'; sillypp 'Hello world!'", &rubyFailed);
  
  if (rubyFailed) {
    //VALUE errMessage = Qnil;
    VALUE errMessage = rb_errinfo();
//    rb_set_errinfo(Qnil);
    
    if (errMessage != Qnil) {
      printf("%s\n", StringValueCStr(errMessage));
    } else {
      printf("No error message\n");
    }
    printf("test failded!\n");
  }
  
  return NULL;
}