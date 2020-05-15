// ANSI-C go<->ruby wrapper 

// see: https://docs.ruby-lang.org/en/2.5.0/extension_rdoc.html
// see: https://ipython.readthedocs.io/en/stable/development/wrapperkernels.html
// see: https://ipython.org/ipython-doc/3/notebook/nbformat.html
// see: https://ipython.org/ipython-doc/dev/development/messaging.html

// see: https://silverhammermba.github.io/emberb

// see: https://gist.github.com/ammar/2787174 for examples of rb_funcall useage

// requires sudo apt install ruby-dev

#include <stdint.h>
#include <ruby.h>
#include <ruby/version.h>
#include "_cgo_export.h"

#include "rubyEval.h"

#ifndef RSTRING_P
#define RSTRING_P(anObj) RB_TYPE_P(anObj, T_STRING)
#endif

/// \brief A global variable which remembers if a Ruby instance has been 
/// started.
///
static int rubyRunning = 0;

/// \brief Starts running the (single) instance of Ruby.
///
void startRuby() {
  // init_ruby....
  if (! rubyRunning) {
    int argc = 0;
    char **argv = 0;
    ruby_sysinit(&argc, &argv);
    RUBY_INIT_STACK;
    ruby_init();
    ruby_init_loadpath();
  
  // load the IPyRubyData.rb file as a string in rb_evalString... IN TOPLEVEL_BINDING
  // need to worry about which binding it is evaluated in... we want our binding
  // so we can (re)use it later....
    rubyRunning = 1;
  }
}

/// \brief Stops running the (single) Ruby instance.
///
int stopRuby() {
  if (rubyRunning) {
    rubyRunning = 0;
    return ruby_cleanup(0);
  }
  return 0;
}

/// \brief Is ruby already running?
///
int isRubyRunning(void) {
  return (int)rubyRunning;
}

#ifndef RUBY_API_VERSION
#define RUBY_API_VERSION \
    STRINGIZE(RUBY_API_VERSION_MAJOR) "." \
    STRINGIZE(RUBY_API_VERSION_MINOR) "." \
    STRINGIZE(RUBY_API_VERSION_TEENY) ""
#endif

/// \brief Returns the string formated ruby version.
///
const char *rubyVersion(void) {
  return RUBY_API_VERSION;
}

/// \brief Create a new Data object and store it in the IPyKernelStore.
///
VALUE IPyKernelData_New(VALUE recv) {
  return LONG2FIX(GoIPyKernelData_New());
}

/// \brief Adds MIMEType/value pair to the Data map of a Data object.
///
///
/// Takes the Data object at `objId` from the IPyKernelStore and adds the 
/// mimeType/dataValue to the Data's Data map.
///
/// The Object ID is a FIXNUM, the mimeType and dataValue are both Ruby 
/// Strings. The dataValue may contain embedded null bytes. 
///
void IPyKernelData_AddData(
  VALUE self,
  VALUE mimeTypeObj,
  VALUE dataValueObj
) {
  uint64_t objId        = 0;
  char    *mimeType     = "";
  uint64_t mimeTypeLen  = 0;
  char    *dataValue    = "";
  uint64_t dataValueLen = 0;

  // check each argument.... do nothing if they are not valid
  //
  if (FIXNUM_P(self)) {
    objId     = NUM2LONG(self);
  } else return ;
  //
  if (RSTRING_P(mimeTypeObj)) {
    mimeType    = StringValuePtr(mimeTypeObj);
    mimeTypeLen = RSTRING_LEN(mimeTypeObj);
  } else return;
  //
  if (RSTRING_P(dataValueObj)) {
    dataValue    = StringValuePtr(dataValueObj);
    dataValueLen = RSTRING_LEN(dataValueObj);
  } else return;

  GoIPyKernelData_AddData(objId, mimeType, mimeTypeLen, dataValue, dataValueLen);
}

/// \brief Adds MIMEType/MetaKey/dataValue triple to the Metadata map of a 
/// Data object. 
///
/// Takes the Data object at `objId` from the IPyKernelStore and adds the 
/// mimeType/metaKey/dataValue to the Data's Metadata map. 
///
/// The Object ID is a FIXNUM, the mimeType, metaKey and dataValue are all 
/// Ruby Strings. The dataValue may contain embedded null bytes. 
///
void IPyKernelData_AddMetadata(
  VALUE self,
  VALUE mimeTypeObj,
  VALUE metaKeyObj,
  VALUE dataValueObj
) {
  uint64_t objId        = 0;
  char    *mimeType     = "";
  uint64_t mimeTypeLen  = 0;
  char    *metaKey      = "";
  uint64_t metaKeyLen   = 0;
  char    *dataValue    = "";
  uint64_t dataValueLen = 0;

  // check each argument.... do nothing if they are not valid
  //
  if (FIXNUM_P(self)) {
    objId     = NUM2LONG(self);
  } else return ;
  //
  if (RSTRING_P(mimeTypeObj)) {
    mimeType    = StringValuePtr(mimeTypeObj); 
    mimeTypeLen = RSTRING_LEN(mimeTypeObj);
  } else return;
  //
  if (RSTRING_P(metaKeyObj)) {
    metaKey    = StringValuePtr(metaKeyObj);
    metaKeyLen = RSTRING_LEN(metaKeyObj);
  } else return;
  //
  if (RSTRING_P(dataValueObj)) {
    dataValue    = StringValuePtr(dataValueObj);
    dataValueLen = RSTRING_LEN(dataValueObj);
  } else return;
  
  GoIPyKernelData_AddMetadata(
    objId,
    mimeType,
    mimeTypeLen,
    metaKey,
    metaKeyLen,
    dataValue,
    dataValueLen
  );
}

/// \brief Evaluate the string aStr in the TOPLEVEL_BINDING and returns 
/// any result as a Go Data object located in the IPyKernelStore at the 
/// returned objId. 
///
uint64_t goIPyRubyEvalString(const char* aStr) {
  //
  // aStr.IPyRubyEval() returns an IPyKernelData object id
  //
  //return rb_funcall(CStr->VALUE(aStr), rb_intern("IPyRubyEval"), 0, Qnil);
  return 0;
}

