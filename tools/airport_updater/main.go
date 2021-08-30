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
	fmt.Println("dataDir:", targetDir)
	alphafoxtrot.DownloadDatabase(targetDir)
}
