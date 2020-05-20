# A Toolkit for building IPython kernels in GoLang

This [IPython kernel](https://ipython.org/) building toolkit is based upon 
the excellent 
[gopherdata/gophernotes](https://github.com/gopherdata/gophernotes), taken 
on 2020-05-08, and used via [gophernotes' MIT 
License](https://github.com/gopherdata/gophernotes/blob/master/LICENSE).

Our goal is to extract gophernotes' embedded Go IPython kernel, and provide 
our kernel with a clean interface to allow multiple implementations.

This project will provide example interface implementations for 
[Lua](https://www.lua.org/), [Ruby](http://www.ruby-lang.org/en/), as well 
as (a heavily refactored) version of Gophernotes itself. This ensures the 
interface we develop is adequate for the job at hand.

## Why not github-fork [gopherdata/gophernotes](https://github.com/gopherdata/gophernotes)?

Why have we forked at the "file level" rather than at the "git level"?

When I started this project I thought carefully about this choice. The 
conclusion that I came to was: 

1. **Gophernotes** is focused upon providing an excellent Jupyter notebook
   kernel for the Go language. This is a clear and simple focus, which 
   *should not be distracted* by trying to track the changes in the 
   goIPythonKernelToolkit project.

2. The **goIPythonKernelToolkit** project is focused upon finding and 
   developing a toolkit and associciated interface for building *multiple* 
   Jupyter notebook kernels for *many different* programming langauges 
   embedded *in* the Go language. Again, this is a clear and simple focus, but 
   a very different focus from Gophernotes' focus. As such, the focus of *this 
   project* should not be distracted by trying to keep up-to-date with the 
   existing Gophernotes project. 

If you are only interested in an excellent Jupyter notebook kernel for the 
Go langauge, *please* use 
[gophernotes](https://github.com/gopherdata/gophernotes). It is, and will 
always be, the better Jupyter notebook kernel for Go. 

If you are interested in building Jupyter notebook kernels for a language 
embedded in Go, then this project might be the best place to start your own 
work.

## External Requirements 

The `kernels/goIPyRuby` and `kernels/goIPyLua` make use of Ruby and Lua 
scripts. These scripts are embedded into the respective go kernel binaries 
using [Matt Jibson's esc tool](https://github.com/mjibson/esc). In order 
to *build* the kernels from source, the `esc` tool must be installed. To 
do this, type: 

```
    go get -u github.com/mjibson/esc
```

Then in each of the respective adaptor directories type:

```
    go generate
```

Then

```
    go install
```

## Notes

We need some way to store the GoLang adaptor or kernel structures for use 
inside the apaptor's ANSI-C code. This would allow state variables to be 
stores on an adaptor or kernel instance scope instead of a global scope.
