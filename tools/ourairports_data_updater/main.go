package main

import (
	"flag"
	"fmt"

	alphafoxtrot "github.com/grumpypixel/go-airport-finder"
)

func main() {
	var targetDir string
	flag.StringVar(&targetDir, "targetdir", "", "OurAirports data directory")
	flag.Parse()
	fmt.Printf("Downloading OurAirports database to: %s\n", targetDir)
	alphafoxtrot.DownloadDatabase(targetDir)
}
