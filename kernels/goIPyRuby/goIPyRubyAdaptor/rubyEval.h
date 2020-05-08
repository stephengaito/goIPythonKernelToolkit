// ANSI-C go<->ruby wrapper (Header)

typedef struct CIPythonReturn {
  const char *mimeType,
  const char *value
} CIPythonReturn;

extern CIPythonReturn *evalString(const char* aStr);

extern freeIPythonReturn(CIPythonRetur *returnVale);