// ANSI-C go<->ruby wrapper 

// see: https://docs.ruby-lang.org/en/2.5.0/extension_rdoc.html
// see: https://ipython.readthedocs.io/en/stable/development/wrapperkernels.html
// see: https://ipython.org/ipython-doc/3/notebook/nbformat.html
// see: https://ipython.org/ipython-doc/dev/development/messaging.html

// see: https://silverhammermba.github.io/emberb

// see: https://gist.github.com/ammar/2787174 for examples of rb_funcall useage

// requires sudo apt install ruby-dev

#include <ruby/ruby.h>
#include <ruby/version.h>
#include "rubyEval.h"

#ifndef RUBY_VERSION
#define RUBY_VERSION \
    STRINGIZE(RUBY_VERSION_MAJOR) "." \
    STRINGIZE(RUBY_VERSION_MINOR) "." \
    STRINGIZE(RUBY_VERSION_TEENY) ""
#endif

const char *rubyVersion(void) {
  return RUBY_VERSION;
}

VALUE IPyKernelData_New(VALUE recv) {
  return LONG2FIX(goIPyKernelData_New());
}

void IPyKernelData_AddData(VALUE self, VALUE mimeType, VALUE dataValue) {
  uint64 objId    = 0;
  char *mimeType  = "";
  char *dataValue = "";

  // check each argument.... do nothing if they are not valid
  //
  if Check_Type(self, T_FIXNUM)      { objId     = NUM2LONG(self);             } else return ;
  if Check_Type(mimeType, T_STRING)  { mimeType  = StringValueCStr(mimeType);  } else return;
  if Check_Type(dataValue, T_STRING) { dataValue = StringValueCStr(dataValue); } else return;

  goIPyKernData_AddData(objId, mimeType, dataValue);
}

void IPyKernelData_AddMetadate(VALUE self, VALUE mimeType, VALUE metaKey, VALUE dataValue) {
  uint64 objId    = 0;
  char *mimeType  = "";
  char *metaKey   = "";
  char *dataValue = "";

  // check each argument.... do nothing if they are not valid
  //
  if Check_Type(self,      T_FIXNUM) { objId     = NUM2LONG(self);             } else return ;
  if Check_Type(mimeType,  T_STRING) { mimeType  = StringValueCStr(mimeType);  } else return;
  if Check_Type(metaKey,   T_STRING) { metaKey   = StringValueCStr(metaKey);   } else return;
  if Check_Type(dataValue, T_STRING) { dataValue = StringValueCStr(dataValue); } else return;
  
  goIPyKernelData_AddMetadata(objId, mimeType, metaKey, dataValue);
}

uint64_t goIPyRubyEvalString(const char* aStr) {
  //
  // aStr.IPyRubyEval() returns an IPyKernelData object id
  //
  return rb_funcall(CStr->VALUE(aStr), rb_intern("IPyRubyEval"), 0, Qnil);
}

