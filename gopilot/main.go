package main

import (
	"app/aeroports"
	"app/filepacker"
	"app/webserver"
	"app/websockets"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/buger/jsonparser"
	"github.com/grumpypixel/msfs2020-simconnect-go/simconnect"
)

type Parameters struct {
	connectionName  string
	searchPath      string
	serverAddress   string
	requestInterval int64
	timeout         int64
}

type Message struct {
	Type  string                 `json:"type"`
	Meta  string                 `json:"meta"`
	Data  map[string]interface{} `json:"data"`
	Debug string                 `json:"debug"`
}

const (
	appTitle                   = "MSFS2020-GoPilot"
	assetsDir                  = "./assets/"
	dataDir                    = "./data/"
	contentTypeHTML            = "text/html"
	contentTypeText            = "text/plain; charset=utf-8"
	defaultServerAddress       = "0.0.0.0:8888"
	defaultSearchPath          = "."
	defaultAirportSearchRadius = 50 * 1000.0
	defaultMaxAirportCount     = 10
	projectURL                 = "http://github.com/grumpypixel/msfs2020-gopilot"
	releasesURL                = projectURL + "/releases"
	connectionTimeout          = 600 // in seconds
	connectRetryInterval       = 1   // in seconds
	requestDataInterval        = 250 // in milliseconds
	receiveDataInterval        = 1   // in milliseconds
	broadcastInterval          = 250
)

type App struct {
	simconnect.EventListener
	requestManager   *RequestManager
	socket           *websockets.WebSocket
	mate             *simconnect.SimMate
	airportsDB       *aeroports.Database
	flightSimVersion string
	done             chan interface{}
}

func main() {
	fmt.Printf("\nWelcome to %s\nProject page: %s\nReleases: %s\n\n", appTitle, projectURL, releasesURL)
	params := &Parameters{}
	parseParameters(params)
	dumpParameters(params)
	checkInstallation(params.searchPath)

	app := NewApp()
	app.run(params)

	fmt.Println("Bye")
}

func parseParameters(params *Parameters) {
	flag.StringVar(&params.connectionName, "name", connectionName(), "Connection name")
	flag.StringVar(&params.searchPath, "searchpath", defaultSearchPath, "Additional DLL search path")
	flag.StringVar(&params.serverAddress, "address", defaultServerAddress, "Web server address (<ipaddr>:<port>)")
	flag.Int64Var(&params.requestInterval, "requestinterval", requestDataInterval, "Request data interval in milliseconds")
	flag.Int64Var(&params.timeout, "timeout", connectionTimeout, "Timeout in seconds")
	flag.Parse()
}

func dumpParameters(params *Parameters) {
	fmt.Println("Application parameters")
	fmt.Println(" Connection name:", params.connectionName)
	fmt.Println(" Additional DLL search path:", params.searchPath)
	fmt.Println(" Web server address:", params.serverAddress)
	fmt.Printf(" RequestInterval: %ds\n", params.requestInterval)
	fmt.Printf(" Timeout: %ds\n\n", params.timeout)
}

func checkInstallation(dllSearchPath string) {
	// Check DLL
	if simconnect.LocateLibrary(dllSearchPath) == false {
		fullpath := path.Join(dllSearchPath, simconnect.SimConnectDLL)
		fmt.Println("DLL not found...")
		data := PackedSimConnectDLL()
		if err := unpack(data, fullpath); err != nil {
			fmt.Println("Unable to unpack DLL:", err)
		}
	}
	// Check assets directory
	if _, err := os.Stat(assetsDir); os.IsNotExist(err) {
		fmt.Println("Assets not found...")
		tarball := "assets.tar"
		fullpath := path.Join("", tarball)
		data := PackedAssets()
		if err := unpack(data, fullpath); err != nil {
			fmt.Println(err)
			return
		}
		if err := filepacker.Untar(tarball, ""); err != nil {
			fmt.Println(err)
			return
		}
		if err := os.Remove(fullpath); err != nil {
			fmt.Println(err)
		}
	}
	// Check data directory
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		fmt.Println("Data not found...")
		tarball := "data.tar"
		fullpath := path.Join("", tarball)
		data := PackedData()
		if err := unpack(data, fullpath); err != nil {
			fmt.Println(err)
			return
		}
		if err := filepacker.Untar(tarball, ""); err != nil {
			fmt.Println(err)
			return
		}
		if err := os.Remove(fullpath); err != nil {
			fmt.Println(err)
		}
	}
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

func NewApp() *App {
	db := aeroports.NewDatabase()
	if err := db.ParseAirports("data/ourairports/airports.csv", aeroports.AirportTypeAll, true); err != nil {
		fmt.Println(err)
		db = nil
	}
	return &App{
		requestManager: NewRequestManager(),
		done:           make(chan interface{}, 1),
		airportsDB:     db,
	}
}

func (app *App) run(params *Parameters) {
	app.socket = websockets.NewWebSocket()
	go app.handleSocketMessages()

	serverShutdown := make(chan bool, 1)
	defer close(serverShutdown)
	app.initWebServer(params.serverAddress, serverShutdown)

	if err := simconnect.Initialize(params.searchPath); err != nil {
		panic(err)
	}

	app.mate = simconnect.NewSimMate()

	stopBroadcast := make(chan interface{}, 1)
	defer close(stopBroadcast)
	go app.Broadcast(time.Millisecond*broadcastInterval, stopBroadcast)

	retryInterval := time.Second * connectRetryInterval
	timeout := time.Second * time.Duration(params.timeout)
	if err := app.connect(params.connectionName, retryInterval, timeout); err != nil {
		fmt.Println(err)
		return
	}

	go app.handleTerminationSignal()

	stopEventHandler := make(chan interface{}, 1)
	defer close(stopEventHandler)

	requestInterval := time.Millisecond * time.Duration(params.requestInterval)
	receiveInterval := time.Millisecond * receiveDataInterval
	go app.mate.HandleEvents(requestInterval, receiveInterval, stopEventHandler, app)

	<-app.done
	defer close(app.done)

	fmt.Println("Shutting down")

	stopBroadcast <- true
	stopEventHandler <- true

	if err := app.disconnect(); err != nil {
		panic(err)
	}
	serverShutdown <- true
}

func (app *App) initWebServer(address string, shutdown chan bool) {
	htmlHeaders := app.Headers(contentTypeHTML)
	textHeaders := app.Headers(contentTypeText)
	webServer := webserver.NewWebServer(address, shutdown)
	routes := []webserver.Route{
		{Pattern: "/", Handler: app.StaticContentHandler(htmlHeaders, "/", filepath.Join("assets/html", "vfrmap.html"))},
		{Pattern: "/vfrmap", Handler: app.StaticContentHandler(htmlHeaders, "/vfrmap", filepath.Join("assets/html", "vfrmap.html"))},
		{Pattern: "/mehmap", Handler: app.StaticContentHandler(htmlHeaders, "/mehmap", filepath.Join("assets/html", "mehmap.html"))},
		{Pattern: "/setdata", Handler: app.StaticContentHandler(htmlHeaders, "/setdata", filepath.Join("assets/html", "setdata.html"))},
		{Pattern: "/airports", Handler: app.StaticContentHandler(htmlHeaders, "/airports", filepath.Join("assets/html", "airports.html"))},
		{Pattern: "/teleport", Handler: app.StaticContentHandler(htmlHeaders, "/teleport", filepath.Join("assets/html", "teleporter.html"))},
		{Pattern: "/debug", Handler: app.GeneratedContentHandler(textHeaders, "/debug", app.DebugGenerator)},
		{Pattern: "/simvars", Handler: app.GeneratedContentHandler(textHeaders, "/simvars", app.SimvarsGenerator)},
		{Pattern: "/ws", Handler: app.socket.Serve},
	}
	fmt.Println("Starting web server")
	assetsDir := "/assets/"
	webServer.Run(routes, assetsDir)

	fmt.Printf("Web Server listening on %s\n\n", address)
	fmt.Println("Your network interfaces:")
	webServer.ListNetworkInterfaces()
	fmt.Println("")
}

func (app *App) connect(name string, retryInterval, timeout time.Duration) error {
	fmt.Println("Connecting to the Simulator...")
	connectTicker := time.NewTicker(retryInterval)
	defer connectTicker.Stop()

	timeoutTimer := time.NewTimer(timeout)
	defer timeoutTimer.Stop()

	count := 0
	for {
		select {
		case <-connectTicker.C:
			count++
			if err := app.mate.Open(name); err != nil {
				if count%10 == 0 {
					fmt.Printf("Connection attempts... %d\n", count)
				}
			} else {
				return nil
			}

		case <-timeoutTimer.C:
			return fmt.Errorf("Opening a connection to the simulator timed out")
		}
	}
}

func (app *App) disconnect() error {
	fmt.Println("Closing connection")
	if err := app.mate.Close(); err != nil {
		return err
	}
	return nil
}

func (app *App) handleTerminationSignal() {
	sigterm := make(chan os.Signal, 1)
	defer close(sigterm)

	signal.Notify(sigterm, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-sigterm:
			fmt.Println("Received SIGTERM")
			app.done <- true
			return
		}
	}
}

func (app *App) handleSocketMessages() {
	for {
		select {
		case event := <-app.socket.EventReceiver:
			eventType := event.Type
			connID := event.Connection.UUID()
			switch eventType {
			case websockets.SocketEventConnected:
				fmt.Println("Client connected:", connID)

			case websockets.SocketEventDisconnected:
				fmt.Println("Client disconnected:", connID)
				app.removeRequests(connID)

			case websockets.SocketEventMessage:
				msg := &Message{}
				json.Unmarshal(event.Data, msg)
				fmt.Println("Message", connID, msg)
				switch msg.Type {
				case "airports":
					app.handleAirportsMessage(msg, connID)

				case "deregister":
					app.handleDeregisterMessage(msg, connID)

				case "echo":
					app.handleEchoMessage(msg, connID)

				case "ping":
					app.handlePingMessage(msg, connID)

				case "register":
					app.handleRegisterMessage(msg, event.Data, connID)

				case "setdata":
					app.handleSetDataMessage(msg)

				case "teleport":
					app.handleTeleportMessage(msg)

				default:
					fmt.Printf("Received unknown message with type: %s\n data: %v\n sender: %s\n", msg.Type, msg.Data, connID)
				}
			}
		}
	}
}

func (app *App) handleAirportsMessage(msg *Message, connID string) {
	if app.airportsDB == nil {
		fmt.Println("Airports database not available")
		return
	}

	latitude, ok := floatFromJson("latitude", msg.Data)
	if !ok {
		return
	}
	longitude, ok := floatFromJson("longitude", msg.Data)
	if !ok {
		return
	}
	radius, ok := floatFromJson("radius", msg.Data)
	if !ok {
		radius = defaultAirportSearchRadius
	}
	maxAirports, ok := intFromJson("maxAirports", msg.Data)
	if !ok {
		maxAirports = defaultMaxAirportCount
	}

	airportFilter := aeroports.AirportTypeAll
	filter, ok := stringFromJson("filter", msg.Data)
	if ok {
		airportFilter = 0
		filters := strings.Split(filter, "|")
		for _, str := range filters {
			f := aeroports.AirportTypeFromString(str)
			airportFilter |= f
		}
	}

	go func() {
		airports := app.airportsDB.FindNearestAirports(latitude, longitude, radius, maxAirports, airportFilter)

		airportList := make([]map[string]interface{}, 0)
		for _, airport := range airports {
			ap := make(map[string]interface{})
			ap["type"] = aeroports.AirportTypeToString(airport.Type)
			ap["icao"] = airport.ICAO
			ap["name"] = airport.Name
			ap["latitude"] = airport.Latitude
			ap["longitude"] = airport.Longitude
			ap["elevation"] = airport.Elevation
			airportList = append(airportList, ap)
		}

		reply := map[string]interface{}{
			"type": "airports",
			"meta": msg.Meta,
			"data": airportList,
		}
		if buf, err := json.Marshal(reply); err == nil {
			app.socket.Send(connID, buf)
		}
	}()
}

func (app *App) handleDeregisterMessage(msg *Message, connID string) {
	app.removeRequests(connID)
}

func (app *App) handleEchoMessage(msg *Message, connID string) {
	if buf, err := json.Marshal(msg); err == nil {
		app.socket.Send(connID, buf)
	}
}

func (app *App) handlePingMessage(msg *Message, connID string) {
	reply := map[string]interface{}{
		"type": "pong",
		"meta": msg.Meta,
		"data": time.Now().String,
	}
	if buf, err := json.Marshal(reply); err == nil {
		app.socket.Send(connID, buf)
	}
}

func (app *App) handleRegisterMessage(msg *Message, raw []byte, connID string) {
	request := NewRequest(connID, msg.Meta)
	jsonparser.ArrayEach(raw, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		n, _, _, _ := jsonparser.Get(value, "name")
		u, _, _, _ := jsonparser.Get(value, "unit")
		t, _, _, _ := jsonparser.Get(value, "type")
		m, _, _, _ := jsonparser.Get(value, "moniker")
		name := string(n)
		unit := string(u)
		typ := simconnect.StringToDataType(string(t))
		moniker := string(m)
		defineID := app.mate.AddSimVar(name, unit, typ)
		request.Add(defineID, name, moniker)
	}, "data")
	app.requestManager.AddRequest(request)
	fmt.Println("Added request", request)
}

func (app *App) handleSetDataMessage(msg *Message) {
	if !app.mate.IsConnected() {
		fmt.Println("Not connected to SimConnect. Ignoring setdata message.")
		return
	}
	name, ok := stringFromJson("name", msg.Data)
	if !ok {
		return
	}
	unit, ok := stringFromJson("unit", msg.Data)
	if !ok {
		return
	}
	value, ok := floatFromJson("value", msg.Data)
	if !ok {
		return
	}
	if err := app.mate.SetSimObjectData(name, unit, value, simconnect.DataTypeFloat64); err != nil {
		fmt.Println(err)
	}
}

func (app *App) handleTeleportMessage(msg *Message) {
	if !app.mate.IsConnected() {
		fmt.Println("Not connected to SimConnect. Ignoring teleport message.")
		return
	}
	latitude, ok := floatFromJson("latitude", msg.Data)
	if !ok {
		return
	}
	longitude, ok := floatFromJson("longitude", msg.Data)
	if !ok {
		return
	}
	altitude, ok := floatFromJson("altitude", msg.Data)
	if !ok {
		return
	}
	heading, ok := floatFromJson("heading", msg.Data)
	if !ok {
		return
	}
	airspeed, ok := floatFromJson("airspeed", msg.Data)
	if !ok {
		return
	}

	bank := 0.0
	pitch := 0.0

	app.mate.SetSimObjectData("PLANE LATITUDE", "degrees", latitude, simconnect.DataTypeFloat64)
	app.mate.SetSimObjectData("PLANE LONGITUDE", "degrees", longitude, simconnect.DataTypeFloat64)
	app.mate.SetSimObjectData("PLANE ALTITUDE", "feet", altitude, simconnect.DataTypeFloat64)
	app.mate.SetSimObjectData("PLANE HEADING DEGREES TRUE", "degrees", heading, simconnect.DataTypeFloat64)
	app.mate.SetSimObjectData("AIRSPEED TRUE", "knot", airspeed, simconnect.DataTypeFloat64)
	app.mate.SetSimObjectData("PLANE BANK DEGREES", "degrees", bank, simconnect.DataTypeFloat64)
	app.mate.SetSimObjectData("PLANE PITCH DEGREES", "degrees", pitch, simconnect.DataTypeFloat64)

	fmt.Printf("Teleporting lat: %f lng: %f alt: %f hdg: %f spd: %f bnk: %f pit: %f\n",
		latitude, longitude, altitude, heading, airspeed, bank, pitch)
}

func (app *App) removeRequests(connID string) {
	temp := make([]*Request, 0)
	removed := make([]*Request, 0)
	for _, request := range app.requestManager.Requests {
		if request.ClientID != connID {
			temp = append(temp, request)
		} else {
			removed = append(removed, request)
		}
	}
	app.requestManager.Requests = temp
	for _, request := range removed {
		for defineID, v := range request.Vars {
			count := app.requestManager.RefCount(v.Name)
			if count == 0 {
				app.mate.RemoveSimVar(defineID)
			}
		}
	}
}

func (app *App) OnOpen(applName, applVersion, applBuild, simConnectVersion, simConnectBuild string) {
	fmt.Println("\nConnected")
	app.flightSimVersion = fmt.Sprintf(
		"Flight Simulator:\n Name: %s\n Version: %s (build %s)\n SimConnect: %s (build %s)",
		applName, applVersion, applBuild, simConnectVersion, simConnectBuild)
	fmt.Printf("\n%s\n\n", app.flightSimVersion)
	fmt.Printf("CLEAR PROP!\n\n")
}

func (app *App) OnQuit() {
	fmt.Println("Disconnected")
	app.done <- true
}

func (app *App) OnEventID(eventID simconnect.DWord) {
	fmt.Println("Received event ID", eventID)
}

func (app *App) OnException(exceptionCode simconnect.DWord) {
	fmt.Printf("Exception (code: %d)\n", exceptionCode)
}

func (app *App) OnDataUpdate(defineID simconnect.DWord, value interface{}) {
	// pass
}

func (app *App) OnDataReady() {
	for _, request := range app.requestManager.Requests {
		msg := map[string]interface{}{
			"type": "simvars",
			"meta": request.Meta,
		}

		vars := make(map[string]interface{})
		for defineID, v := range request.Vars {
			value, dataType, ok := app.mate.SimVarValueAndDataType(defineID)
			if !ok || value == nil {
				continue
			}
			switch dataType {
			case simconnect.DataTypeInt32:
				vars[v.Moniker] = simconnect.ValueToInt32(value)

			case simconnect.DataTypeInt64:
				vars[v.Moniker] = simconnect.ValueToInt64(value)

			case simconnect.DataTypeFloat32:
				vars[v.Moniker] = simconnect.ValueToFloat32(value)

			case simconnect.DataTypeFloat64:
				vars[v.Moniker] = simconnect.ValueToFloat64(value)

			case simconnect.DataTypeString8,
				simconnect.DataTypeString32,
				simconnect.DataTypeString64,
				simconnect.DataTypeString128,
				simconnect.DataTypeString256,
				simconnect.DataTypeString260,
				simconnect.DataTypeStringV:
				vars[v.Moniker] = simconnect.ValueToString(value)
			}
		}
		msg["data"] = vars
		recipient := request.ClientID
		if buf, err := json.Marshal(msg); err == nil {
			app.socket.Send(recipient, buf)
		}
	}
}

func (app *App) Broadcast(broadcastInterval time.Duration, stop chan interface{}) {
	broadcastTicker := time.NewTicker(broadcastInterval)
	defer broadcastTicker.Stop()

	for {
		select {
		case <-broadcastTicker.C:
			if err := app.BroadcastStatusMessage(); err != nil {
				fmt.Println(err)
			}
		case <-stop:
			return
		}
	}
}

func (app *App) BroadcastStatusMessage() error {
	data := map[string]interface{}{"simconnect": app.mate.IsConnected()}
	msg := map[string]interface{}{"type": "status", "data": data}
	buf, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	app.socket.Broadcast(buf)
	return nil
}

func (app *App) Headers(contentType string) map[string]string {
	headers := map[string]string{
		"Access-Control-Allow-Origin": "*",
		"Cache-Control":               "no-cache, no-store, must-revalidate",
		"Pragma":                      "no-cache",
		"Expires":                     "0",
		"Content-Type":                contentType,
	}
	return headers
}

func (app *App) StaticContentHandler(headers map[string]string, urlPath, filePath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != urlPath {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		for key, value := range headers {
			w.Header().Set(key, value)
		}
		http.ServeFile(w, r, filePath)
	}
}

func (app *App) GeneratedContentHandler(headers map[string]string, urlPath string, generator func(w http.ResponseWriter)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != urlPath {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		for key, value := range headers {
			w.Header().Set(key, value)
		}
		w.WriteHeader(http.StatusOK)
		generator(w)
	}
}

func (app *App) SimvarsGenerator(w http.ResponseWriter) {
	fmt.Fprintf(w, "%s\n\n", appTitle)
	fmt.Fprintf(w, "%s\n", app.DumpedSimVars())
}

func (app *App) DebugGenerator(w http.ResponseWriter) {
	fmt.Fprintf(w, "%s\n\n", appTitle)
	if len(app.flightSimVersion) > 0 {
		fmt.Fprintf(w, "%s\n\n", app.flightSimVersion)
	}
	fmt.Fprintf(w, "SimConnect\n  initialized: %v\n  conncected: %v\n\n", simconnect.IsInitialized(), app.mate.IsConnected())
	fmt.Fprintf(w, "Clients: %d\n", app.socket.ConnectionCount())
	uuids := app.socket.ConnectionUUIDs()
	for i, uuid := range uuids {
		fmt.Fprintf(w, "  %02d: %s\n", i, uuid)
	}
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "%s\n\n", app.DumpedSimVars())
	fmt.Fprintf(w, "%s\n", app.DumpedRequests())
}

func (app *App) DumpedRequests() string {
	var dump string
	dump += fmt.Sprintf("Requests: %d\n", app.requestManager.RequestCount())
	for i, request := range app.requestManager.Requests {
		dump += fmt.Sprintf("  %02d: Client: %s Vars: %d Meta: %s\n", i+1, request.ClientID, len(request.Vars), request.Meta)
		count := 1
		for name, moniker := range request.Vars {
			dump += fmt.Sprintf("    %02d: name: %s moniker: %s\n", count, name, moniker)
			count++
		}
	}
	return dump
}

func (app *App) DumpedSimVars() string {
	indent := "  "
	dump := app.mate.SimVarDump(indent)
	str := strings.Join(dump[:], "\n")
	return fmt.Sprintf("SimVars: %d\n", len(dump)) + str
}

func floatFromJson(key string, json map[string]interface{}) (float64, bool) {
	value, ok := json[key]
	if !ok {
		return 0.0, false
	}
	return value.(float64), true
}

func intFromJson(key string, json map[string]interface{}) (int, bool) {
	value, ok := json[key]
	if !ok {
		return 0, false
	}
	return int(value.(float64)), true
}

func stringFromJson(key string, json map[string]interface{}) (string, bool) {
	value, ok := json[key]
	if !ok {
		return "", false
	}
	return value.(string), true
}

// func (app *App) handleSetCameraMessage(msg *Message) {
// 	fmt.Println("SetCameraMessage", *msg)
// 	deltaX := msg.Data["delta_x"].(float64)
// 	deltaY := msg.Data["delta_y"].(float64)
// 	deltaZ := msg.Data["delta_z"].(float64)
// 	pitch := msg.Data["pitch"].(float64)
// 	bank := msg.Data["bank"].(float64)
// 	heading := msg.Data["heading"].(float64)
// 	app.mate.CameraSetRelative6DOF(deltaX, deltaY, deltaZ, pitch, bank, heading)
// }

// func (app *App) handleSetTextMessage(msg *Message) {
// 	fmt.Println("SetTextMessage", *msg)
// 	// IMPLEMENT ME
// 	text := "HELLO, SIMWORLD!"
// 	textType := simconnect.TextTypePrintMagenta
// 	duration := 10.0
// 	eventID := simconnect.NextEventID()
// 	app.mate.Text(text, textType, duration, eventID)
// }
