// ANSI-C go<->ruby wrapper (Header)

typedef struct CIPythonReturn {
  const char *mimeType;
  const char *value;
} CIPythonReturn;

extern const char *rubyVersion(void);

extern CIPythonReturn *evalString(const char* aStr);

extern void freeIPythonReturn(CIPythonReturn *returnValue);