# goIPythonKernelToolkit documentation

This is a complex project. It is *critical* that we expose as much of the 
structure of the code as possible, preferably with good (internal) 
documentation.

This directory 


## Requirements

We use Hugo to manage the local Markdown file rendering. Since you must 
already have GoLang installed to use this project, to add Hugo type: 

```
    go get github.com/gohugoio/hugo
```
(Or grab a binary release from Hugo)

The `docs/runDocServers` script uses the `webfsd` static webserver to add 
it (on a Debian based OS) type:

```
    sudo apt install webfs
```


