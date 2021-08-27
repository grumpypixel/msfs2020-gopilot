package main

import (
	"flag"
	"fmt"
	"math/rand"
	"msfs2020-gopilot/internal/app"
	"msfs2020-gopilot/internal/config"
	"msfs2020-gopilot/internal/filepacker"
	"os"
	"path"
	"strings"
	"time"

	"github.com/grumpypixel/msfs2020-simconnect-go/simconnect"
	"github.com/mattn/go-colorable"
	log "github.com/sirupsen/logrus"
)

const (
	appTitle                   = "MSFS2020-GoPilot"
	assetsDir                  = "assets/"
	dataDir                    = "data/"
	contentTypeHTML            = "text/html"
	contentTypeText            = "text/plain; charset=utf-8"
	defaultServerAddress       = "0.0.0.0:8888"
	defaultSearchPath          = "."
	defaultAirportSearchRadius = 50 * 1000.0
	defaultMaxAirportCount     = 10
	projectURL                 = "http://github.com/grumpypixel/msfs2020-gopilot"
	releasesURL                = projectURL + "/releases"
	connectionTimeout          = 600 // seconds
	connectRetryInterval       = 1   // seconds
	requestDataInterval        = 250 // milliseconds
	receiveDataInterval        = 1   // milliseconds
	shutdownDelay              = 3   // seconds
	broadcastInterval          = 250
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		ForceColors: true,
		// DisableColors: false,
		// FullTimestamp: false,
	})
	log.SetOutput(colorable.NewColorableStdout())
}

func main() {
	fmt.Printf("\nWelcome to %s\nVisit: %s\nReleases: %s\n\n", appTitle, projectURL, releasesURL)
	params := &config.Parameters{}
	parseParameters(params)
	dumpParameters(params)
	if err := checkInstallation(params.SearchPath); err != nil {
		log.Fatal(err)
	}

	app := app.NewApp()
	app.Run(params)

	fmt.Println("Bye")
}

func connectionName() string {
	rand.Seed(time.Now().Unix())
	names := []string{
		"0xDECAFBAD", "0xBADDCAFE", "0xCAFED00D",
		"Boobytrap", "Sobeit Void", "Transpotato",
		"A but Tuba", "Evil Olive", "Flee to Me, Remote Elf",
		"Sit on a Potato Pan, Otis", "Taco Cat", "UFO Tofu",
	}
	return names[rand.Intn(len(names))]
}

func parseParameters(params *config.Parameters) {
	flag.StringVar(&params.ConnectionName, "name", connectionName(), "Connection name")
	flag.StringVar(&params.SearchPath, "searchpath", defaultSearchPath, "Additional DLL search path")
	flag.StringVar(&params.ServerAddress, "address", defaultServerAddress, "Web server address (<ipaddr>:<port>)")
	flag.Int64Var(&params.RequestInterval, "requestinterval", requestDataInterval, "Request data interval in milliseconds")
	flag.Int64Var(&params.Timeout, "timeout", connectionTimeout, "Timeout in seconds")
	// boolean params expect an equal sign (=) between the variable name and the value, i.e. verbose=true. meh.
	// see also: https://stackoverflow.com/questions/27411691/how-to-pass-boolean-arguments-to-go-flags/27411724
	// flag.BoolVar(&params.verbose, "verbose", false, "Verbosity")
	// so out of pure convenience we'll use strings here
	verbose := flag.String("verbose", "false", "Verbosity")
	flag.Parse()

	*verbose = strings.ToLower(*verbose)
	params.Verbose = *verbose == "1" || *verbose == "true"
}

func dumpParameters(params *config.Parameters) {
	fmt.Println("Application parameters")
	fmt.Println(" Connection name:", params.ConnectionName)
	fmt.Println(" Additional DLL search path:", params.SearchPath)
	fmt.Println(" Web server address:", params.ServerAddress)
	fmt.Printf(" RequestInterval: %ds\n", params.RequestInterval)
	fmt.Printf(" Timeout: %ds\n", params.Timeout)
	fmt.Printf(" Verbosity: %v\n\n", params.Verbose)
}

func checkInstallation(dllSearchPath string) error {
	// Check DLL
	if !simconnect.LocateLibrary(dllSearchPath) {
		fullpath := path.Join(dllSearchPath, simconnect.SimConnectDLL)
		fmt.Println("DLL not found...")
		data := app.PackedSimConnectDLL()
		if err := unpack(data, fullpath); err != nil {
			log.Error("Unable to unpack DLL:", err)
		}
	}
	// Check assets directory
	if _, err := os.Stat(assetsDir); os.IsNotExist(err) {
		return err
	}
	// Check data directory
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		return err
	}
	return nil
}

func unpack(data []byte, fullpath string) error {
	fmt.Printf("Unpacking target: %s\n", fullpath)
	unpacked, err := filepacker.Unpack(data)
	if err != nil {
		return err
	}
	file, err := os.Create(fullpath)
	if err != nil {
		return err
	}
	if _, err := file.WriteString(string(unpacked)); err != nil {
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	time.Sleep(time.Second * 1)
	return nil
}
