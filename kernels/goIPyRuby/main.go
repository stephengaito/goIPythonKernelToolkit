package main

import (
	"flag"
	"log"
  
  "github.com/stephengaito/goIPythonKernelToolkit/goIPyKernel"
  "github.com/stephengaito/goIPythonKernelToolkit/kernels/goIPyRuby/goIPyRubyAdaptor"
)

func main() {

	// Parse the connection file.
	flag.Parse()
	if flag.NArg() < 1 {
		log.Fatalln("Need a command line argument specifying the connection file.")
	}

  adaptor := goIPyRubyAdaptor.NewGoAdaptor()
  kernel  := goIPyKernel.NewIPyKernel(adaptor)
  
	// Run the kernel.
	kernel.Run(flag.Arg(0))
}