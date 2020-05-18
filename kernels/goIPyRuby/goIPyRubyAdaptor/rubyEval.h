// ANSI-C go<->ruby wrapper (Header)

#ifndef RUBY_EVAL_H
#define RUBY_EVAL_H

extern void startRuby(void);

extern int stopRuby(void);

extern int isRubyRunning(void);

typedef struct LoadRubyCodeReturn_struct {
  char     *errMesg;
  uint64_t  objId;
} LoadRubyCodeReturn;

LoadRubyCodeReturn *FreeLoadRubyCodeReturn(LoadRubyCodeReturn *aReturn);

extern LoadRubyCodeReturn *loadRubyCode(
  const char *rubyCodeNameCStr,
  const char *rubyCodeCStr
);

extern int isRubyCodeLoaded(
  const char *rubyCodeNameCStr
);

extern const char *rubyVersion(void);

extern uint64_t evalRubyString(
  const char* evalNameCStr,
  const char* evalCodeCStr
);

#endif
