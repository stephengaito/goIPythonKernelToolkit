// ANSI-C go<->ruby wrapper 

// see: https://docs.ruby-lang.org/en/2.5.0/extension_rdoc.html
// see: https://ipython.readthedocs.io/en/stable/development/wrapperkernels.html
// see: https://ipython.org/ipython-doc/3/notebook/nbformat.html
// see: https://ipython.org/ipython-doc/dev/development/messaging.html

// see: https://silverhammermba.github.io/emberb

// see: https://gist.github.com/ammar/2787174 for examples of rb_funcall useage

// assertions: https://godoc.org/github.com/stretchr/testify/assert

// uthash: https://troydhanson.github.io/uthash/userguide.html

// requires sudo apt install ruby-dev

#include <assert.h>
#include <stdint.h>
#include <pthread.h>
#include <ruby.h>
#include <ruby/version.h>
#include "_cgo_export.h"

#include "rubyEval.h"
#define HASH_FUNCTION HASH_FNV
#include "uthash.h"

#define DEBUG
#ifdef DEBUG
#define DEBUG_Log(aMessage) printf(aMessage); fflush(stdout)
#define DEBUG_Log2(aFormat, aValue) printf(aFormat, aValue); fflush(stdout)
#define DEBUG_Log3(aFormat, aValue, anotherValue)             \
  printf(aFormat, aValue, anotherValue); fflush(stdout)
#else
#define DEBUG_Log(aMessage) 
#define DEBUG_Log2(aFormat, aValue)
#define DEBUG_Log3(aFormat, aValue, anotherValue)
#endif

void Init_IPyKernelData(void); // forward declaration...

#ifndef RSTRING_P
#define RSTRING_P(anObj) RB_TYPE_P(anObj, T_STRING)
#endif

/// \brief A global variable which remembers if a Ruby instance has been 
/// started.
///
static int rubyRunning = 0;

/// \brief We use a PThreads mutex to ensure only one thread uses Ruby at 
/// any one time. 
///
static pthread_mutex_t rubyMutex   = PTHREAD_MUTEX_INITIALIZER;


/// \brief the uthash structure to hold the hash table (entries). 
///
typedef struct LoadedCodeNames {
  const char     *codeName;
  UT_hash_handle  hh;
} LoadedCodeNames;

/// \brief We use uthash to keep a hash table of all already loaded
/// RubyCodeNames. This ensures we do not get multiply defined symbol 
/// errors. 
///
static LoadedCodeNames *loadedCodeNames = NULL;

/// \brief Starts running the (single) instance of Ruby.
///
void startRuby() {
  // init_ruby....
  pthread_mutex_lock(&rubyMutex);
  if (! rubyRunning) {
    printf("Starting Ruby\n");
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
  pthread_mutex_unlock(&rubyMutex);
}

/// \brief Stops running the (single) Ruby instance.
///
int stopRuby() {
  pthread_mutex_lock(&rubyMutex);
  if (rubyRunning) {
    rubyRunning = 0;
    return ruby_cleanup(0);
  }
  pthread_mutex_unlock(&rubyMutex);
  return 0;
}

/// \brief Is ruby already running?
///
int isRubyRunning(void) {
  return (int)rubyRunning;
}

LoadRubyCodeReturn *FreeLoadRubyCodeReturn(LoadRubyCodeReturn *aReturn) {
  if (aReturn->errMesg) {
    free(aReturn->errMesg);
  }
  free(aReturn);
  return NULL;
}

static VALUE protectedLoadRubyCode(VALUE args) {
  VALUE rubyCodeStr  = rb_ary_pop(args);
  VALUE rubyCodeName = rb_ary_pop(args);
  
  DEBUG_Log("before protectedLoadRubyCode::rb_funcall\n");
  VALUE result = rb_funcall(
    Qnil,
    rb_intern("eval"),
    4,
    rubyCodeStr,
    Qnil,
    rubyCodeName,
    LONG2FIX(0),
    0
  );
  // This will NOT be called IF the eval raises an exception...
  DEBUG_Log("after protectedLoadRubyCode::rb_funcall\n");
  return result;
}

/// \brief Load the Ruby code from the string provided
///
LoadRubyCodeReturn *loadRubyCode(
  const char *rubyCodeNameCStr,
  const char *rubyCodeCStr
) {
  DEBUG_Log("START loadRubyCode\n");
  DEBUG_Log2(" rubyCodeName: %s\n", rubyCodeNameCStr);
  //DEBUG_Log2("     rubyCode: %s\n", rubyCodeCStr);

  pthread_mutex_lock(&rubyMutex);

  
  LoadRubyCodeReturn *returnStruct = calloc(1, sizeof(LoadRubyCodeReturn));
  DEBUG_Log2("          returnStruct: %p\n", returnStruct);
  DEBUG_Log2("   returnStruct->objId: %ld\n", returnStruct->objId);
  DEBUG_Log2(" returnStruct->errMesg: %s\n", returnStruct->errMesg);

  LoadedCodeNames *foundCodeName = NULL;
  DEBUG_Log3(
    "Looking for [%s] in uthash %p in loadRubyCode\n",
    rubyCodeNameCStr,
    loadedCodeNames
  );
  HASH_FIND_STR(loadedCodeNames, rubyCodeNameCStr, foundCodeName);
  if (!foundCodeName) {
    //
    // this code has not yet been loaded... so load it...
    //
    DEBUG_Log2("Need to load: %s\n", rubyCodeNameCStr);
    VALUE rubyCodeName = rb_str_new_cstr(rubyCodeNameCStr);
    VALUE rubyCodeStr  = rb_str_new_cstr(rubyCodeCStr);
    VALUE codeArray = rb_ary_new();
    DEBUG_Log2("rubyCodeName: %ld\n", rubyCodeName);
    DEBUG_Log2(" rubyCodeStr: %ld\n", rubyCodeStr);
    DEBUG_Log2("   codeArray: %ld\n", codeArray);
    rb_ary_push(codeArray, rubyCodeName);
    rb_ary_push(codeArray, rubyCodeStr);
  
    DEBUG_Log("Before rb_protect\n");
    int loadFailed = 0;
    VALUE result = rb_protect(protectedLoadRubyCode, codeArray, &loadFailed);
    DEBUG_Log2("After rb_protect     result: %ld\n", result);
    DEBUG_Log2("After rb_protect loadFailed: %d\n", loadFailed);
    if (loadFailed) {
      DEBUG_Log("Load failed\n");
      VALUE errMesg = rb_errinfo();
      DEBUG_Log2("errMesg: %ld\n", errMesg);
      VALUE errStr  = rb_sprintf("%"PRIsVALUE, errMesg);
      DEBUG_Log2("errStr: %ld\n", errStr);
      DEBUG_Log2("%s", StringValueCStr(errStr));
      returnStruct->errMesg = 
        strndup(StringValuePtr(errStr), RSTRING_LEN(errStr));
    } else {
      DEBUG_Log3(
        "adding [%s] to uthash %p in loadRubyCode\n", 
        rubyCodeNameCStr, 
        loadedCodeNames
      );
      LoadedCodeNames *newCodeName = calloc(1, sizeof(LoadedCodeNames));
      assert(newCodeName);
      newCodeName->codeName = strdup(rubyCodeNameCStr);
      HASH_ADD_STR(loadedCodeNames, codeName, newCodeName);
    }
    DEBUG_Log2("result %ld\n", result);
    if (RB_FIXNUM_P(result)) {
      returnStruct->objId = FIX2LONG(result);
    }
    DEBUG_Log2("Finished loading: %s\n", rubyCodeNameCStr);
  } else {
    DEBUG_Log2("Found [%s] in loadRubyCode\n", rubyCodeNameCStr);
  }
  pthread_mutex_unlock(&rubyMutex);

  DEBUG_Log("FINISHED loadRubyCode\n");
  return returnStruct;
}

int isRubyCodeLoaded(const char *rubyCodeNameCStr) {
  int result = 0;
  
  pthread_mutex_lock(&rubyMutex);

  LoadedCodeNames *foundCodeName = NULL;
  DEBUG_Log3(
    "Looking for [%s] in uthash %p in isRubyCodeLoaded\n",
    rubyCodeNameCStr,
    loadedCodeNames
  );
  HASH_FIND_STR(loadedCodeNames, rubyCodeNameCStr, foundCodeName);
  if (foundCodeName) { 
    DEBUG_Log2("Found [%s] in isRubyCodeLoaded\n", rubyCodeNameCStr);
    result = 1;
  }

  pthread_mutex_unlock(&rubyMutex);

  return result ;
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
  DEBUG_Log2("  objId %ld\n", newObjId);
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
  DEBUG_Log2("  objId: %ld\n", objId);
  //
  if (RSTRING_P(mimeTypeObj)) {
    mimeType    = StringValuePtr(mimeTypeObj);
    mimeTypeLen = RSTRING_LEN(mimeTypeObj);
  } else return Qnil;
  DEBUG_Log2("  mimeType: %s\n", mimeType);
  //
  if (RSTRING_P(dataValueObj)) {
    dataValue    = StringValuePtr(dataValueObj);
    dataValueLen = RSTRING_LEN(dataValueObj);
  } else return Qnil;
  DEBUG_Log2("  dataValue: %s\n", dataValue);

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
  DEBUG_Log2("  objId: %ld\n", objId);
  //
  if (RSTRING_P(tracebackValueObj)) {
    tracebackValue    = StringValuePtr(tracebackValueObj);
    tracebackValueLen = RSTRING_LEN(tracebackValueObj);
  } else return Qnil;
  DEBUG_Log2("  tracebackValue: %s\n", tracebackValue);

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
  DEBUG_Log2("  objId: %ld\n", objId);
  //
  if (RSTRING_P(mimeTypeObj)) {
    mimeType    = StringValuePtr(mimeTypeObj); 
    mimeTypeLen = RSTRING_LEN(mimeTypeObj);
  } else return Qnil;
  DEBUG_Log2("  mimeType: %s\n", mimeType);
  //
  if (RSTRING_P(metaKeyObj)) {
    metaKey    = StringValuePtr(metaKeyObj);
    metaKeyLen = RSTRING_LEN(metaKeyObj);
  } else return Qnil;
  DEBUG_Log2("  metaKey: %s\n", metaKey);
  //
  if (RSTRING_P(dataValueObj)) {
    dataValue    = StringValuePtr(dataValueObj);
    dataValueLen = RSTRING_LEN(dataValueObj);
  } else return Qnil;
  DEBUG_Log2("  dataValue: %s\n", dataValue);
  
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
  DEBUG_Log("before protectedEvalString::rb_funcall\n");
  VALUE result = rb_funcall(
    Qnil,
    rb_intern("IPyRubyEval"),
    4,
    evalStr,
    Qnil,
    evalName,
    LONG2FIX(0),
    0
  );
  // This will NOT be called IF the IPyRubyEval raises an exception...
  DEBUG_Log("after protectedEvalString::rb_funcall\n");
  return result;
}

/// \brief Evaluate the string aStr in the TOPLEVEL_BINDING and returns 
/// any result as a Go Data object located in the IPyKernelStore at the 
/// returned objId. 
///
uint64_t evalRubyString(
  const char* evalNameCStr,
  const char* evalCodeCStr
) {
  DEBUG_Log2("Starting evalRubyString on [%s]\n", evalNameCStr);
  pthread_mutex_lock(&rubyMutex);

  VALUE evalName = rb_str_new_cstr(evalNameCStr);
  VALUE evalCode = rb_str_new_cstr(evalCodeCStr);
  
  VALUE evalArray = rb_ary_new();
  rb_ary_push(evalArray, evalName);
  rb_ary_push(evalArray, evalCode);
  
  DEBUG_Log("Before rb_protect\n");
  int loadFailed = 0;
  uint64_t result = 0;
  VALUE rbResult = rb_protect(protectedEvalString, evalArray, &loadFailed);
  if (RB_FIXNUM_P(rbResult)) { result = FIX2LONG(rbResult); }
  DEBUG_Log2("After rb_protect   rbResult: %ld\n", rbResult);
  DEBUG_Log2("After rb_protect     result: %ld\n", result);
  DEBUG_Log2("After rb_protect loadFailed: %d\n", loadFailed);  
  if (loadFailed) {
    GoIPyKernelData_Delete(result);

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
  
  pthread_mutex_unlock(&rubyMutex);
  DEBUG_Log2("Finished evalRubyString on [%s]\n", evalNameCStr);
  return result;
}

