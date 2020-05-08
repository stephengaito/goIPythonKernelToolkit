package goIPyLuaAdaptor


// #include "luaEval.h"
import "C"

type LuaState uintptr

func newLuaInterpreter() LuaState {
  return C.newLuaInterpreter()
}

func (l LuaState) closeLuaInterpreter() {
  C.closeLuaInterpreter(l);
}

func (l LuaState) evalString(aGoStr string) {
  const char* aCStr = C.CString(aGoStr)
  defer C.free(aCStr)

  C.evalString(l, aCStr)
}
