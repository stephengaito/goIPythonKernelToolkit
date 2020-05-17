// ANSI-C go<->ruby wrapper (Header)

extern void startRuby(void);

extern int stopRuby(void);

extern int isRubyRunning(void);

extern char *loadRubyCode(
  const char *rubyCodeNameCStr,
  const char *rubyCodeCStr
);

extern int isRubyCodeLoaded(
  const char *rubyCodeNameCStr
);

extern const char *rubyVersion(void);

extern uint64_t evalString(const char* aStr);
