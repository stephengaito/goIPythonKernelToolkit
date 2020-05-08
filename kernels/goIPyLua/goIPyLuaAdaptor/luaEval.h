// ANSI-C go<->lua wrapper (Header)

uintptr_t newLuaInterpreter(void);
void closeLuaInterpreter(uintptr_t theInterp);
void evalString(uintptr_t theInterp, const char* aStr);
