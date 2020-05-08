// ANSI-C go<->lua wrapper

// requires sudo apt install lua5.x-dev

#include <stdio.h>
#include <string.h>
#include "lua.h"
#include "lauxlib.h"
#include "lualib.h"

type LuaState struct {
  Interp *C.lua_State
}

uintptr_t newLuaInterpreter(void) {
  lua_State *L = luaL_newstate();
  luaL_openlibs(L);
  return (uintptr_t)L;
}

void closeLuaInterpreter(uintptr_t theInterp) {
  lua_close((lua_State*)theInterp);
}

void evalString(uintptr_t theInterp, const char* aStr) {
  lua_State *L = (lua_State*)theInterp;
  error =
    luaL_loadstring(L, aStr) ||
    lua_pcall(L, 0, LUA_MULTRET, 0);
    
// using LUA_MULTRET ensures that any "returned" values are pushed onto 
// the stack (to be popped off below). (See documentation for lua_pcall 
// and lua_call) 

  if (error) {
    fprintf(stderr, "%s\n", C.lua_tostring(l.Interp, -1));
    C.lua_pop(l.Interp, 1); /* pop error message from the stack */
  }
  
// Then use lua_gettop to determine how many values to pop off the lua 
// stack and return to jupyter... for each value use the lua_type method 
// to determine what type the current top of the stack represents and then 
// get that type as needed. 

// Assemble the reulsts into an array.... and then use a Lua script to 
// encode the result into either a 'pp' or json string. 

}

