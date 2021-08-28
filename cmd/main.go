package main

import (
	"encoding/json"
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

	"github.com/common-nighthawk/go-figure"
	"github.com/grumpypixel/msfs2020-simconnect-go/simconnect"
	"github.com/mattn/go-colorable"
	log "github.com/sirupsen/logrus"
)

const (
	appTitle  = "MSFS2020-GoPilot"
	dataDir   = "data/"
	assetsDir = "assets/"

	defaultConfigFilePath = "./config/config.yaml"
	defaultServerAddress  = "0.0.0.0:8888"
	defaultDLLSearchPath  = "."
	projectURL            = "http://github.com/grumpypixel/msfs2020-gopilot"
	releasesURL           = projectURL + "/releases"
	connectionTimeout     = 600 // seconds
	requestDataInterval   = 200 // milliseconds
)

type Parameters struct {
	ConfigFilePath      string
	ConnectionName      string
	ConnectionTimeout   int64
	DLLSearchPath       string
	ServerAddress       string
	DataRequestInterval int64
	Verbose             bool
}

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})
	log.SetOutput(colorable.NewColorableStdout())
}

func main() {
	welcome()

	params := parseParams()
	prettyPrint("Your command line parameters:\n", params)

	cfg, err := config.NewConfig(params.ConfigFilePath)
	if err != nil {
		log.Fatal(err)
	}
	prettyPrint("Your configuration:\n", cfg)

	mergeConfig(params, cfg)
	validateConfig(cfg)
	prettyPrint("Final configuration:\n", cfg)

	if err := checkInstallation(params.DLLSearchPath); err != nil {
		log.Fatal(err)
	}

	app := app.NewApp(cfg)
	if err := app.Run(); err != nil {
		log.Error(err)
	}

	log.Info("Bye \\o/")
}

func welcome() {
	asciiLogo := figure.NewFigure("GoPilot", "doom", true)
	asciiLogo.Print()
	fmt.Printf("\nWelcome to %s\nHomepage: %s\nReleases: %s\n\n", appTitle, projectURL, releasesURL)
}

func parseParams() *Parameters {
	params := &Parameters{}
	flag.StringVar(&params.ConfigFilePath, "config", defaultConfigFilePath, "Config file location")
	flag.StringVar(&params.ConnectionName, "name", "", "Connection name")
	flag.StringVar(&params.DLLSearchPath, "searchpath", "", "Additional DLL search path")
	flag.StringVar(&params.ServerAddress, "address", "", "Web server address (ipaddr:port)")
	flag.Int64Var(&params.DataRequestInterval, "requestinterval", -1, "Request data interval in milliseconds")
	flag.Int64Var(&params.ConnectionTimeout, "timeout", -1, "Timeout in seconds")

	// boolean params expect an equal sign (=) between the variable name and the value, i.e. verbose=true. meh.
	// see also: https://stackoverflow.com/questions/27411691/how-to-pass-boolean-arguments-to-go-flags/27411724
	// flag.BoolVar(&params.verbose, "verbose", false, "Verbosity")
	// so out of pure convenience we'll use strings here
	verbose := flag.String("verbose", "false", "Verbosity")
	flag.Parse()

	*verbose = strings.ToLower(*verbose)
	params.Verbose = *verbose == "1" || *verbose == "true"
	return params
}

func mergeConfig(params *Parameters, cfg *config.Config) {
	if params.ConnectionName != "" {
		cfg.ConnectionName = params.ConnectionName
	}
	if params.ConnectionTimeout >= 0 {
		cfg.ConnectionTimeout = params.ConnectionTimeout
	}
	if params.DLLSearchPath != "" {
		cfg.DLLSearchPath = params.DLLSearchPath
	}
	if params.ServerAddress != "" {
		cfg.ServerAddress = params.ServerAddress
	}
	if params.DataRequestInterval >= 0 {
		cfg.DataRequestInterval = params.DataRequestInterval
	}
	if params.Verbose {
		cfg.Verbose = params.Verbose
	}
}

func validateConfig(cfg *config.Config) {
	if cfg.ConnectionName == "" {
		cfg.ConnectionName = randomConnectionName()
	}
}

func randomConnectionName() string {
	rand.Seed(time.Now().Unix())
	names := []string{
		"0xDECAFBAD", "0xBADDCAFE", "0xCAFED00D",
		"Boobytrap", "Sobeit Void", "Transpotato",
		"A but Tuba", "Evil Olive", "Flee to Me, Remote Elf",
		"Sit on a Potato Pan, Otis", "Taco Cat", "UFO Tofu",
	}
	return names[rand.Intn(len(names))]
}

func checkInstallation(dllSearchPath string) error {
	// Check DLL
	if !simconnect.LocateLibrary(dllSearchPath) {
		fullpath := path.Join(dllSearchPath, simconnect.SimConnectDLL)
		log.Warn("DLL not found...")
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

func prettyPrint(info string, data interface{}) {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Error(err)
		return
	}
	log.Info(info, string(bytes))
}
