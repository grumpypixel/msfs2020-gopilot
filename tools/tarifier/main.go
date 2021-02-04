package main

import (
	"app/filepacker"
	"flag"
)

func main() {
	var input, output string
	flag.StringVar(&input, "in", "", "Input directory")
	flag.StringVar(&output, "out", "", "Output file")
	flag.Parse()

	if err := filepacker.Tar(input, output); err != nil {
		panic(err)
	}
}
