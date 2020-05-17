// ANSI-C go<->ruby wrapper 

// see: https://docs.ruby-lang.org/en/2.5.0/extension_rdoc.html
// see: https://ipython.readthedocs.io/en/stable/development/wrapperkernels.html
// see: https://ipython.org/ipython-doc/3/notebook/nbformat.html
// see: https://ipython.org/ipython-doc/dev/development/messaging.html

// see: https://silverhammermba.github.io/emberb

// see: https://gist.github.com/ammar/2787174 for examples of rb_funcall useage

// assertions: https://godoc.org/github.com/stretchr/testify/assert

// requires sudo apt install ruby-dev

#include <stdint.h>
#include <ruby.h>
#include <ruby/version.h>
#include "_cgo_export.h"

#include "rubyEval.h"

//#define DEBUG_Log(aMessage) printf(aMessage)
//#define DEBUG_Logf(aFormat, aValue) printf(aFormat, aValue)
#define DEBUG_Log(aMessage) 
#define DEBUG_Logf(aFormat, aValue)

void Init_IPyKernelData(void); // forward declaration...

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
    ruby_script("goIPyRuby");
    Init_IPyKernelData();
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

static VALUE protectedLoadRubyCode(VALUE args) {
  VALUE rubyCodeStr  = rb_ary_pop(args);
  VALUE rubyCodeName = rb_ary_pop(args);
  
  return rb_funcall(
    Qnil,
    rb_intern("eval"),
    4,
    rubyCodeStr,
    Qnil,
    rubyCodeName,
    LONG2FIX(0),
    0
  );
}

VALUE loadedCodeNames = Qnil; 

/// \brief Load the Ruby code from the string provided
///
char *loadRubyCode(const char *rubyCodeNameCStr, const char *rubyCodeCStr) {
  if (loadedCodeNames == Qnil) { loadedCodeNames = rb_hash_new(); }

  VALUE rubyCodeStr  = rb_str_new_cstr(rubyCodeCStr);
  VALUE rubyCodeName = rb_str_new_cstr(rubyCodeNameCStr);
  VALUE isLoaded     = rb_hash_aref(loadedCodeNames, rubyCodeName);
  if (isLoaded == Qtrue) { return 0; }
  
  VALUE codeArray = rb_ary_new();
  rb_ary_push(codeArray, rubyCodeName);
  rb_ary_push(codeArray, rubyCodeStr);
  
  int loadFailed = 0;
  VALUE result = rb_protect(protectedLoadRubyCode, codeArray, &loadFailed);
  if (loadFailed) { 
    VALUE errMesg = rb_errinfo();
    VALUE errStr  = rb_sprintf("%"PRIsVALUE, errMesg);
    
    DEBUG_Logf("%s", StringValueCStr(errStr));
    return strndup(StringValuePtr(errStr), RSTRING_LEN(errStr));
  }
  rb_hash_aset(loadedCodeNames, rubyCodeName, Qtrue);
  return 0;
}

int isRubyCodeLoaded(const char *rubyCodeNameCStr) {
  if (loadedCodeNames == Qnil) { loadedCodeNames = rb_hash_new(); }

  VALUE rubyCodeName = rb_str_new_cstr(rubyCodeNameCStr);
  VALUE isLoaded     = rb_hash_aref(loadedCodeNames, rubyCodeName);
  if (isLoaded == Qtrue) { return 1; }
  return 0;
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
  DEBUG_Log("IPyKernelData_New\n");
  uint64_t newObjId = GoIPyKernelData_New();
  DEBUG_Logf("  objId %ld\n", newObjId);
  return  LONG2FIX(newObjId);
}

/// \brief Adds MIMEType/value pair to the Data map of a Data object.
///
/// Takes the Data object at `objId` from the IPyKernelStore and adds the 
/// mimeType/dataValue to the Data's Data map.
///
/// The Object ID is a FIXNUM, the mimeType and dataValue are both Ruby 
/// Strings. The dataValue may contain embedded null bytes. 
///
VALUE IPyKernelData_AddData(
  VALUE recv,
  VALUE objIdObj,
  VALUE mimeTypeObj,
  VALUE dataValueObj
) {
  uint64_t objId        = 0;
  char    *mimeType     = "";
  uint64_t mimeTypeLen  = 0;
  char    *dataValue    = "";
  uint64_t dataValueLen = 0;

  DEBUG_Log("IPyKernelData_AddData\n");
  // check each argument.... do nothing if they are not valid
  //
  if (FIXNUM_P(objIdObj)) {
    objId     = NUM2LONG(objIdObj);
  } else return Qnil;
  DEBUG_Logf("  objId: %ld\n", objId);
  //
  if (RSTRING_P(mimeTypeObj)) {
    mimeType    = StringValuePtr(mimeTypeObj);
    mimeTypeLen = RSTRING_LEN(mimeTypeObj);
  } else return Qnil;
  DEBUG_Logf("  mimeType: %s\n", mimeType);
  //
  if (RSTRING_P(dataValueObj)) {
    dataValue    = StringValuePtr(dataValueObj);
    dataValueLen = RSTRING_LEN(dataValueObj);
  } else return Qnil;
  DEBUG_Logf("  dataValue: %s\n", dataValue);

  GoIPyKernelData_AddData(objId, mimeType, mimeTypeLen, dataValue, dataValueLen);
  return Qnil;
}

/// \brief Appends a traceback error message to the `traceback` entry in 
/// the Data map of an (error) Data object. 
///
/// Takes the Data object at `objId` from the IPyKernelStore and appends the 
/// tracebackValue to the Data's `traceback` entry in the Data map.
///
/// The Object ID is a FIXNUM, the tracebackValue is a Ruby 
/// Strings.
///
VALUE IPyKernelData_AppendTraceback(
  VALUE recv,
  VALUE objIdObj,
  VALUE tracebackValueObj
) {
  uint64_t objId              = 0;
  char    *tracebackValue     = "";
  uint64_t tracebackValueLen  = 0;
  DEBUG_Log("IPyKernelData_AppendTraceback\n");
  // check each argument.... do nothing if they are not valid
  //
  if (FIXNUM_P(objIdObj)) {
    objId     = NUM2LONG(objIdObj);
  } else return Qnil;
  DEBUG_Logf("  objId: %ld\n", objId);
  //
  if (RSTRING_P(tracebackValueObj)) {
    tracebackValue    = StringValuePtr(tracebackValueObj);
    tracebackValueLen = RSTRING_LEN(tracebackValueObj);
  } else return Qnil;
  DEBUG_Logf("  tracebackValue: %s\n", tracebackValue);

  GoIPyKernelData_AppendTraceback(objId, tracebackValue, tracebackValueLen);
  return Qnil;
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
VALUE IPyKernelData_AddMetadata(
  VALUE recv,
  VALUE objIdObj,
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

  DEBUG_Log("IPyKernelData_AddMetadata\n");
  // check each argument.... do nothing if they are not valid
  //
  if (FIXNUM_P(objIdObj)) {
    objId     = NUM2LONG(objIdObj);
  } else return Qnil;
  DEBUG_Logf("  objId: %ld\n", objId);
  //
  if (RSTRING_P(mimeTypeObj)) {
    mimeType    = StringValuePtr(mimeTypeObj); 
    mimeTypeLen = RSTRING_LEN(mimeTypeObj);
  } else return Qnil;
  DEBUG_Logf("  mimeType: %s\n", mimeType);
  //
  if (RSTRING_P(metaKeyObj)) {
    metaKey    = StringValuePtr(metaKeyObj);
    metaKeyLen = RSTRING_LEN(metaKeyObj);
  } else return Qnil;
  DEBUG_Logf("  metaKey: %s\n", metaKey);
  //
  if (RSTRING_P(dataValueObj)) {
    dataValue    = StringValuePtr(dataValueObj);
    dataValueLen = RSTRING_LEN(dataValueObj);
  } else return Qnil;
  DEBUG_Logf("  dataValue: %s\n", dataValue);
  
  GoIPyKernelData_AddMetadata(
    objId,
    mimeType,
    mimeTypeLen,
    metaKey,
    metaKeyLen,
    dataValue,
    dataValueLen
  );
  return Qnil;
}

/// \brief Initialize the IPyKernelData class inside ruby
///
void Init_IPyKernelData(void) {
  rb_define_global_function("IPyKernelData_New",             IPyKernelData_New,             1);
  rb_define_global_function("IPyKernelData_AddData",         IPyKernelData_AddData,         3);
  rb_define_global_function("IPyKernelData_AppendTraceback", IPyKernelData_AppendTraceback, 2);
  rb_define_global_function("IPyKernelData_AddMetadata",     IPyKernelData_AddMetadata,     4);
}

static VALUE protectedEvalString(VALUE args) {
  VALUE evalStr  = rb_ary_pop(args);
  VALUE evalName = rb_ary_pop(args);
  
  return rb_funcall(
    Qnil,
    rb_intern("eval"),
    4,
    evalStr,
    Qnil,
    evalName,
    LONG2FIX(0),
    0
  );
}
/// \brief Evaluate the string aStr in the TOPLEVEL_BINDING and returns 
/// any result as a Go Data object located in the IPyKernelStore at the 
/// returned objId. 
///
uint64_t goIPyRubyEvalString(const char* evalCStr, const char* evalNameCStr) {

  VALUE evalStr  = rb_str_new_cstr(evalCStr);
  VALUE evalName = rb_str_new_cstr(evalNameCStr);
  
  VALUE evalArray = rb_ary_new();
  rb_ary_push(evalArray, evalName);
  rb_ary_push(evalArray, evalStr);
  
  int loadFailed = 0;
  uint64_t result = FIX2LONG(rb_protect(protectedLoadRubyCode, evalArray, &loadFailed));
  if (loadFailed) { 
    VALUE errMesg = rb_errinfo();
    VALUE errStr  = rb_sprintf("%"PRIsVALUE, errMesg);
    result = GoIPyKernelData_New();
    GoIPyKernelData_AddData(result,
      "ename", strlen("ename"), "ERROR", strlen("ERROR"));
    GoIPyKernelData_AddData(result,
      "evalue", strlen("evalue"), StringValuePtr(errStr), RSTRING_LEN(errStr));
    char* tracebackMsg = "protectedLoadRubyCode FAILED";
    GoIPyKernelData_AppendTraceback(result,
      tracebackMsg, strlen(tracebackMsg));
    GoIPyKernelData_AddData(result,
      "status", strlen("status"), "error", strlen("status"));
  }
  return result;
}

