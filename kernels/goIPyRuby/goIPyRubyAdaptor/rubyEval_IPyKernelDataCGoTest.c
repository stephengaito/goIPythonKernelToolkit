// +buid cGoTests

// Some tests of the IPyKernelData interface using cGoTests

#include "_cgo_export.h"
#include "goIPyRubyAdaptorCGoTests.h"
#include "cGoTests.h"

void addMimeMapToDataObjTest(uint64_t objId) {
  char *mimeType  = "MIMETest";
  char *dataValue = "some data";
  
  GoIPyKernelData_AddData(
    objId,
    mimeType,  strlen(mimeType),
    dataValue, strlen(dataValue)
  );
}

void addJPEGMimeMapToDataObjTest(uint64_t objId) {
  char *mimeType  = "image/jpeg";
  char dataValue[10];
  dataValue[0] = 's';
  dataValue[1] = 'o';
  dataValue[2] = 'm';
  dataValue[3] = 'e';
  dataValue[4] = 0;
  dataValue[5] = 'd';
  dataValue[6] = 'a';
  dataValue[7] = 't';
  dataValue[8] = 'a';
  dataValue[9] = 0;
  
  GoIPyKernelData_AddData(
    objId,
    mimeType,  strlen(mimeType),
    dataValue, 10
  );
}

void addPNGMimeMapToDataObjTest(uint64_t objId) {
  char *mimeType  = "image/png";
  char dataValue[10];
  dataValue[0] = 's';
  dataValue[1] = 'o';
  dataValue[2] = 'm';
  dataValue[3] = 'e';
  dataValue[4] = 0;
  dataValue[5] = 'd';
  dataValue[6] = 'a';
  dataValue[7] = 't';
  dataValue[8] = 'a';
  dataValue[9] = 0;
  
  GoIPyKernelData_AddData(
    objId,
    mimeType,  strlen(mimeType),
    dataValue, 10
  );
}

void addMimeMapToMetadataObjTest(uint64_t objId) {
  char *mimeType  = "MIMETest";
  char *metaKey   = "Width";
  char *dataValue = "some data";
  
  GoIPyKernelData_AddMetadata(
    objId,
    mimeType,  strlen(mimeType),
    metaKey,   strlen(metaKey),
    dataValue, strlen(dataValue)
  );
}

/// \brief Test something
///
char *IPyKernelDataCGoTest(void* data) {
  uint64_t objId = GoIPyKernelData_New();
  cGoTest_UIntEquals(
    "objId should be TheObjectStore.NextId",
    objId,
    1
  );
  return NULL;
}