<head><title>Refactoring Plan</title></head>

# Refactoring plan

## Display code

Reviewing the internal structure of the gophernotes *.ipynb files suggests 
that **over the wire** only the following mimetypes are of any importance: 

-	MIMETypeHTML       = "text/html"
-	MIMETypeJavaScript = "application/javascript"
-	MIMETypeJPEG       = "image/jpeg"
-	MIMETypeJSON       = "application/json"
-	MIMETypeLatex      = "text/latex"
-	MIMETypeMarkdown   = "text/markdown"
-	MIMETypePNG        = "image/png"
-	MIMETypePDF        = "application/pdf"
-	MIMETypeSVG        = "image/svg+xml"
-	MIMETypeText       = "text/plain"

This suggests that the bulk of the existing Gophernotes "display" code, 
which is used to translate *from* a reflected GoLang type, should be 
moved to the goIPyGophernotes kernel. 

This means that each of the kernels is responsible for translating its 
native types into an appropriate IPython mimetype. For the langauges which 
are embedded via cgo, this will typically take the form of an ANSI-C 
sprintf of a basic ANSI-C type into a string. 

The more sophisticated embedded kernels might use a simple json-printer to 
translate complex structures into a form of json. 

## Handler code

The other code which should be moved from the existing Gophernotes code 
into the goIPyGophernotes code is the code to handle the high level 
message generation (doEval) and formating:

- `sendKernelInfo`
- `handelExecuteRequest`
    - `doEval`
- `handleCompleteRequest` (or some delegate calls)
- `evalSpecialCommands` (`evalSpecialCommand`) (in particular the 
  `%go111module {on|off}` command) 
- `evalShellCommand` (should this remain generic?)

We must also figure out how to handle "exceptions" to be handed back via 
the `executionErr` code in `handleExecuteRequest`. 

## Improving the (internal) documentation

All structure fields and functions *should* begin with capital letters to 
ensure, as "public" parts of the various interfaces, they appear in the 
godoc documentation. 

All structure fields and functions *should* get some documentation (if 
they do not alreay have any). 

