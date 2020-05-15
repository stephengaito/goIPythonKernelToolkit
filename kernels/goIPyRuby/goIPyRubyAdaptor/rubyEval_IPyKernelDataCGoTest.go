// +build cGoTests

package goIPyRubyAdaptor

// #include <stdint.h>
// extern void addMimeMapToDataObjTest(uint64_t objId);
// extern void addJPEGMimeMapToDataObjTest(uint64_t objId);
// extern void addPNGMimeMapToDataObjTest(uint64_t objId);
// extern void addMimeMapToMetadataObjTest(uint64_t objId);
import "C"

func GoAddMimeMapToDataObjTest(objId uint64) {
  C.addMimeMapToDataObjTest(C.ulong(objId));
}

func GoAddJPEGMimeMapToDataObjTest(objId uint64) {
  C.addJPEGMimeMapToDataObjTest(C.ulong(objId));
}

func GoAddPNGMimeMapToDataObjTest(objId uint64) {
  C.addPNGMimeMapToDataObjTest(C.ulong(objId));
}

func GoAddMimeMapToMetadataObjTest(objId uint64) {
  C.addMimeMapToMetadataObjTest(C.ulong(objId));
}
