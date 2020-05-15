// ANSI-C go<->ruby wrapper (Header)

extern void startRuby(void);

extern int stopRuby(void);

extern int isRubyRunning(void);

extern const char *rubyVersion(void);

extern uint64_t evalString(const char* aStr);
