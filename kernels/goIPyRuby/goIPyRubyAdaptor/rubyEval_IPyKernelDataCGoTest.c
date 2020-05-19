// +buid cGoTests

/// \file
/// \brief Some tests of the IPyKernelData interface using cGoTests

#include "_cgo_export.h"
#include "goIPyRubyAdaptorCGoTests.h"
#include "cGoTests.h"


/// \brief Test adding a MIMEType/dataValue to the Data field of a 
/// DataObj.
///
/// For use by the Go rubyEval_IPyKernelDataCGTest.go and
/// rubyEval_IPyKernelData_test.go files. We use this inderection since Go 
/// tests can not directly call CGo wrapped code, and we want to ensure 
/// the ANSI-C code *can* call back into the Go code. 
///
void addMimeMapToDataObjTest(uint64_t objId) {
  char *mimeType  = "MIMETest";
  char *dataValue = "some data";
  
  GoIPyKernelData_AddData(
    objId,
    mimeType,  strlen(mimeType),
    dataValue, strlen(dataValue)
  );
}

/// \brief Test adding a MIMETypeJPEF/dataValue to the Data field of a 
/// DataObj.
///
/// For use by the Go rubyEval_IPyKernelDataCGTest.go and
/// rubyEval_IPyKernelData_test.go files. We use this inderection since Go 
/// tests can not directly call CGo wrapped code, and we want to ensure 
/// the ANSI-C code *can* call back into the Go code. 
///
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

/// \brief Test adding a MIMETypePNG/dataValue to the Data field of a 
/// DataObj.
///
/// For use by the Go rubyEval_IPyKernelDataCGTest.go and
/// rubyEval_IPyKernelData_test.go files. We use this inderection since Go 
/// tests can not directly call CGo wrapped code, and we want to ensure 
/// the ANSI-C code *can* call back into the Go code. 

///
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

/// \brief Test adding a MIMEType/metaKey/dataValue to the Metadata field 
/// of a DataObj.
///
/// For use by the Go rubyEval_IPyKernelDataCGTest.go and
/// rubyEval_IPyKernelData_test.go files. We use this inderection since Go 
/// tests can not directly call CGo wrapped code, and we want to ensure 
/// the ANSI-C code *can* call back into the Go code. 
/// 
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
