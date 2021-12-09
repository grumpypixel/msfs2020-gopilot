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
	"runtime/debug"
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

	defaultConfigFilePath      = "./configs/config.yml"
	defaultServerAddress       = "0.0.0.0:8888"
	defaultSimConnectDLLPath   = "."
	defaultConnectionTimeout   = 600 // seconds
	defaultRequestDataInterval = 200 // milliseconds
	projectURL                 = "http://github.com/grumpypixel/msfs2020-gopilot"
	releasesURL                = projectURL + "/releases"
)

type Parameters struct {
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

	defer func() {
		if r := recover(); r != nil {
			log.Warnf("stacktrace from panic: %s", string(debug.Stack()))
		}
	}()

	var configFilePath string
	flag.StringVar(&configFilePath, "cfg", defaultConfigFilePath, "Config file location")
	flag.Parse()

	log.Infof("Loading config at {%s}", configFilePath)
	cfg, err := config.NewConfigFromFile(configFilePath)
	if err != nil {
		log.Info("Loading a default configuration...")
		cfg = newDefaultConfig()
	}
	prettyPrint("Configuration:\n", cfg)

	log.SetLevel(getLogLevel(cfg.LogLevel))

	if err := checkInstallation(cfg.SimConnectDLLPath); err != nil {
		log.Fatal(err)
	}

	app := app.NewApp(cfg)
	if err := app.Run(); err != nil {
		log.Error(err)
	}

	log.Info("Bye \\(^-^)/")
}

func welcome() {
	asciiLogo := figure.NewFigure("GoPilot", "doom", true)
	asciiLogo.Print()
	fmt.Printf("\nWelcome to %s\nHomepage: %s\nReleases: %s\n\n", appTitle, projectURL, releasesURL)
}

func newDefaultConfig() *config.Config {
	return &config.Config{
		ConnectionName:      randomConnectionName(),
		ConnectionTimeout:   defaultConnectionTimeout,
		SimConnectDLLPath:   ".",
		ServerAddress:       defaultServerAddress,
		DataRequestInterval: defaultRequestDataInterval,
		LogLevel:            "info",
	}
}

func checkInstallation(simConnectDLLPath string) error {
	// Check DLL
	libPath, err := simconnect.LocateLibrary(simConnectDLLPath)
	if libPath != "" {
		log.Infof("Found DLL at {%s}", libPath)
	}
	if err != nil {
		log.Errorf("DLL not found at {%s} (error: %s)", simConnectDLLPath, err.Error())
		fullpath := path.Join("./", simconnect.SimConnectDLL)
		log.Info("Unpacking SimConnect.dll to ", fullpath)
		data := app.SimConnectDLL()
		if err := unpack(data, fullpath); err != nil {
			log.Error("Unable to unpack DLL: ", err)
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

func prettyPrint(msg string, data interface{}) {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Error(err)
		return
	}
	log.Info(msg, string(bytes))
}

func getLogLevel(level string) log.Level {
	switch level {
	case "error":
		return log.ErrorLevel
	case "warn":
		return log.WarnLevel
	case "debug":
		return log.DebugLevel
	case "trace":
		return log.TraceLevel
	}
	return log.InfoLevel
}

func randomConnectionName() string {
	rand.Seed(time.Now().Unix())
	names := []string{
		"0xDECAFBAD", "0xBADDCAFE", "0xCAFED00D",
		"Boobytrap", "Sobeit Void", "Transpotato",
		"A Butt Tuba", "Evil Olive", "Flee to Me, Remote Elf",
		"Sit on a Potato Pan, Otis", "Taco Cat", "UFO Tofu",
	}
	return names[rand.Intn(len(names))]
}
