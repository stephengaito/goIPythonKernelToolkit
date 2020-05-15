package goIPyKernel

import (
  "sync"
)

// A global object store to permit ANSI-C code to interact with long lived 
// Go objects without explicitly keeping Go pointers in C-objects and or 
// C-code. 
//
// Objects in the store are indexed by a uint64 value.
//
type ObjectStore struct {
  Mutex   sync.RWMutex
  NextId  uint64
  Objects map[uint64]interface{}
}

// A global object store to permit ANSI-C code to interact with long lived 
// Go objects without explicitly keeping Go pointers in C-objects and or 
// C-code. 
//
// Objects in the store are indexed by a uint64 value.
//
var TheObjectStore = &ObjectStore{
  Objects: make(map[uint64]interface{}),
}

// Store a new object `aValue` into an ObjectStore.
//
// Panics if there are more than ^uint64(0) - 10 objects.
//
func (theStore *ObjectStore) Store(aValue interface{} ) uint64 {
  theStore.Mutex.Lock()
  defer theStore.Mutex.Unlock()
  
  if theStore.NextId < ( ^uint64(0) - 10 ) {
    theStore.NextId = theStore.NextId + 1
    theStore.Objects[theStore.NextId] = aValue
    return theStore.NextId
  }
  panic("Run out of object ids in the ObjectStore!")
}

func (theStore *ObjectStore) Get(anObjId uint64) interface{} {
  theStore.Mutex.RLock()
  defer theStore.Mutex.RUnlock()
  
  return theStore.Objects[anObjId]
}

func (theStore *ObjectStore) Delete(anObjId uint64) {
  theStore.Mutex.Lock()
  defer theStore.Mutex.Unlock()
  
  if theStore.Objects[anObjId] != nil {
    delete(theStore.Objects, anObjId)
  }
}
