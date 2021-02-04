package main

import (
	"app/filepacker"
	"flag"
	"fmt"
	"time"
)

func main() {
	fmt.Println("Packifying...")
	var infile, outfile, template, packageName, getterName string
	flag.StringVar(&infile, "in", "", "Input file")
	flag.StringVar(&outfile, "out", "", "Output file")
	flag.StringVar(&template, "template", "", "Name of the template file")
	flag.StringVar(&packageName, "package", "main", "Name of the package")
	flag.StringVar(&getterName, "function", "GetData", "Name of the getter")
	flag.Parse()

	if len(infile) == 0 {
		panic("No input file specified")
	}
	if len(outfile) == 0 {
		panic("No input file specified")
	}
	if len(template) == 0 {
		panic("No template file specified")
	}

	fmt.Println(" Input file:", infile)
	fmt.Println(" Ouput file:", outfile)
	fmt.Println(" Template:", template)
	fmt.Println(" Package:", packageName)
	fmt.Println(" Getter:", getterName)

	now := time.Now()
	timestamp := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d.%d",
		now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Second(), now.Nanosecond())
	fmt.Println(" Timestamp:", timestamp)

	filepacker.Pack(infile, outfile, template, timestamp, packageName, getterName)
}
