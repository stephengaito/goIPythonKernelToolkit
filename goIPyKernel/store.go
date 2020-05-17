package goIPyKernel

import (
  "sync"
)

// A lockable object (as stored in the ObjectStore)
//
type LockableObject struct {
  Mutex sync.Mutex
  obj   interface{}
}

// A global object store to permit ANSI-C code to interact with long lived 
// Go objects without explicitly keeping Go pointers in C-objects and or 
// C-code. 
//
// Objects in the store are indexed by a uint64 value.
//
type ObjectStore struct {
  Mutex   sync.RWMutex
  NextId  uint64
  Objects map[uint64]*LockableObject
}

// A global object store to permit ANSI-C code to interact with long lived 
// Go objects without explicitly keeping Go pointers in C-objects and or 
// C-code. 
//
// Objects in the store are indexed by a uint64 value.
//
var TheObjectStore = &ObjectStore{
  Objects: make(map[uint64]*LockableObject),
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
    theStore.Objects[theStore.NextId] = &LockableObject{
      obj: aValue,
    }
    return theStore.NextId
  }
  panic("Run out of object ids in the ObjectStore!")
}

// GetLocked returns the object with the object id of `anObjId` ONLY once 
// that object can be (globally) locked. This call will block until the 
// requested object in NOT (globally) locked. 
//
func (theStore *ObjectStore) GetLocked(anObjId uint64) interface{} {
  if anObjId == 0 { return nil }
  
  theStore.Mutex.RLock()
  theLockableObj := theStore.Objects[anObjId]
  theStore.Mutex.RUnlock()
  
  if theLockableObj == nil { return nil }
  
  theLockableObj.Mutex.Lock()
  return theLockableObj.obj
}

// Unlock unlocks the (global) lock on the object in `theStore` with the 
// object id of `anObjId`. 
//
// NOTE: Panics if there is no (global) lock on the object with object id 
// `anOjbId` in `theStore`. 
//
func (theStore *ObjectStore) Unlock(anObjId uint64) {
  if anObjId == 0 { return }
  
  theStore.Mutex.RLock()
  theLockableObj := theStore.Objects[anObjId]
  theStore.Mutex.RUnlock()
  
  if theLockableObj == nil { return }
  
  theLockableObj.Mutex.Unlock()
}

// NOTE: It is critical that the object being deleted is Unlocked... 
//
func (theStore *ObjectStore) Delete(anObjId uint64) {
  theStore.Mutex.Lock()
  defer theStore.Mutex.Unlock()
  
  if theStore.Objects[anObjId] != nil {
    delete(theStore.Objects, anObjId)
  }
}
