package goIPyKernel

// assertions: 

import(
  "testing"
  "github.com/stretchr/testify/assert"
)

// assertions: https://godoc.org/github.com/stretchr/testify/assert

func TestTheObjectStore(t *testing.T) {
  assert.NotNil(t, TheObjectStore, "TheObjectStore is nil")
  
  // store an object
  anObjId := TheObjectStore.Store("this is a test")
  assert.NotZero(t, anObjId, "ObjId is zero")
  assert.Equal(t, anObjId, uint64(1), "ObjId is not one")
  
  // get that object
  anObj := TheObjectStore.Get(anObjId)
  
  assert.NotNil(t, anObj, "anObj is nil")
  aStr := anObj.(string)
  assert.Equal(t, aStr, "this is a test", "aStr is not as stored")
 
  // show that we can delete an object
  TheObjectStore.Delete(anObjId)

  // show that that object no longer exists
  anObj = TheObjectStore.Get(anObjId)
  assert.Nil(t, anObj, "anObj is not nil")
  
  // show that we can delete twice 
  TheObjectStore.Delete(anObjId)
}
