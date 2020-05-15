package goIPyKernel

import(
  "testing"
)

func TestTheObjectStore(t *testing.T) {
  if TheObjectStore == nil { t.Error("TheObjectStore is nil") }
  
  // store an object
  anObjId := TheObjectStore.Store("this is a test")
  if anObjId == 0 { t.Error("ObjId is zero")}
  if anObjId != 1 { t.Error("ObjId is not one")}
  
  // get that object
  anObj := TheObjectStore.Get(anObjId)
  if anObj == nil { t.Error("anObj is nil")}
  aStr := anObj.(string)
  if aStr != "this is a test" { t.Error("aStr is not as stored")}
 
  // show that we can delete an object
  TheObjectStore.Delete(anObjId)

  // show that that object no longer exists
  anObj = TheObjectStore.Get(anObjId)
  if anObj != nil { t.Error("anObj is not nil")}
  
  // show that we can delete twice 
  TheObjectStore.Delete(anObjId)
  
}
