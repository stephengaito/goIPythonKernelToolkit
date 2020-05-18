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
// NOTE: Panics if there are more than ^uint64(0) - 10 objects.
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

// Get returns the object with the object id of `anObjId`. 
//
func (theStore *ObjectStore) Get(anObjId uint64) interface{} {
  if anObjId == 0 { return nil }
  
  theStore.Mutex.RLock()
  theObj := theStore.Objects[anObjId]
  theStore.Mutex.RUnlock()
  
  return theObj
}

// Delete the object with `anObjId` from the object store.
//
func (theStore *ObjectStore) Delete(anObjId uint64) {
  theStore.Mutex.Lock()
  defer theStore.Mutex.Unlock()
  
  if theStore.Objects[anObjId] != nil {
    delete(theStore.Objects, anObjId)
  }
}
