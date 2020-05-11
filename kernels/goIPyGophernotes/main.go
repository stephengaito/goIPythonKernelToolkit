package main

import (
	"flag"
	"log"
)

const (

	// Version defines the goIPyGophernotes version.
	Version string = "1.0.0"

	// ProtocolVersion defines the Jupyter protocol version.
	ProtocolVersion string = "5.0"
)

func main() {

	// Parse the connection file.
	flag.Parse()
	if flag.NArg() < 1 {
		log.Fatalln("Need a command line argument specifying the connection file.")
	}

  adaptor := NewGoAdaptor()
  kernel  := NewIPyKernel(adaptor)
  
	// Run the kernel.
	kernel.Run(flag.Arg(0))
}
