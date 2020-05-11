// ANSI-C go<->ruby wrapper 

// see: https://docs.ruby-lang.org/en/2.5.0/extension_rdoc.html
// see: https://ipython.readthedocs.io/en/stable/development/wrapperkernels.html
// see: https://ipython.org/ipython-doc/3/notebook/nbformat.html
// see: https://ipython.org/ipython-doc/dev/development/messaging.html

// requires sudo apt install ruby-dev

#include <ruby/ruby.h>
#include "rubyEval.h"

// USE the json.rb library (included in Ruby 2.5.x and up) to encode the 
// results. 
// OR USE the pp.rb library (included in Ruby 2.5.x and up) to encode the
// retsults.

// IF a result is a hash whose keys are string names of the IPython 
// Mimetypes, then the result is turned into the corresponding Go/IPython 
// Mimemap. If a result is not a hash or is a hash whose keys are not
// IPython Mimetypes, then the result is turned into a MIMETypeText using
// the 'pp' script mentioned above.

// Image data can be stored as Ruby strings with embedded zeros, but MUST 
// be returned as the valid MIMETypeJPEG, or MIMETypePNG.

void evalString(const char* aStr) {
  VALUE result = rb_eval_string(aStr);
  switch (TYPE(result)) {
    case T_NIL      : // nil
    case T_OBJECT   : // ordinary object
    case T_CLASS    : // class
    case T_MODULE   : // module
    case T_FLOAT    : // floating point number
    case T_STRING   : // string
    case T_REGEXP   : // regular expression
    case T_ARRAY    : // array
    case T_HASH     : // associative array
    case T_STRUCT   : // (Ruby) structure
    case T_BIGNUM   : // multi precision integer
    case T_FIXNUM   : // Fixnum(31bit or 63bit integer)
    case T_COMPLEX  : // complex number
    case T_RATIONAL : // rational number
    case T_FILE     : // IO
    case T_TRUE     : // true
    case T_FALSE    : // false
    case T_DATA     : // (user) data
    case T_SYMBOL   : // symbol
    case T_ICLASS   : // (internal) included module
    case T_MATCH    : // (internal) MatchData object
    case T_UNDEF    : // (internal) undefined
    case T_NODE     : // (internal) syntax tree node
    case T_ZOMBIE   : // (internal) object awaiting finalization
    default         :
    
  }
}

