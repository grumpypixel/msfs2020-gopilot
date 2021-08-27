package app

import (
	"encoding/json"
	"fmt"
	"msfs2020-gopilot/internal/aeroports"
	"msfs2020-gopilot/internal/config"
	"msfs2020-gopilot/internal/webserver"
	"msfs2020-gopilot/internal/websockets"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/buger/jsonparser"
	"github.com/grumpypixel/msfs2020-simconnect-go/simconnect"
	log "github.com/sirupsen/logrus"
)

type Message struct {
	Type  string                 `json:"type"`
	Meta  string                 `json:"meta"`
	Data  map[string]interface{} `json:"data"`
	Debug string                 `json:"debug"`
}

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

type App struct {
	requestManager   *RequestManager
	socket           *websockets.WebSocket
	mate             *simconnect.SimMate
	airportsDB       *aeroports.Database
	done             chan interface{}
	flightSimVersion string
	verbose          bool
	eventListener    *simconnect.EventListener
}

func NewApp() *App {
	return &App{
		requestManager: NewRequestManager(),
		done:           make(chan interface{}, 1),
		airportsDB:     aeroports.NewDatabase(),
	}
}

func (app *App) Run(params *config.Parameters) {
	app.verbose = params.Verbose

	app.addEventListeners()

	if err := app.airportsDB.ParseAirports(dataDir+"ourairports/airports.csv", aeroports.AirportTypeAll, true); err != nil {
		log.Error(err)
		app.airportsDB = nil
	}

	app.socket = websockets.NewWebSocket()
	go app.handleSocketMessages()

	serverShutdown := make(chan bool, 1)
	defer close(serverShutdown)
	app.initWebServer(params.ServerAddress, serverShutdown)

	if err := simconnect.Initialize(params.SearchPath); err != nil {
		log.Fatal(err)
	}

	app.mate = simconnect.NewSimMate()

	stopBroadcast := make(chan interface{}, 1)
	defer close(stopBroadcast)
	go app.Broadcast(time.Millisecond*broadcastInterval, stopBroadcast)

	retryInterval := time.Second * connectRetryInterval
	timeout := time.Second * time.Duration(params.Timeout)
	if err := app.connect(params.ConnectionName, retryInterval, timeout); err != nil {
		log.Error(err)
		return
	}

	go app.handleTerminationSignal()

	stopEventHandler := make(chan interface{}, 1)
	defer close(stopEventHandler)

	requestInterval := time.Millisecond * time.Duration(params.RequestInterval)
	receiveInterval := time.Millisecond * receiveDataInterval
	go app.mate.HandleEvents(requestInterval, receiveInterval, stopEventHandler, app.eventListener)

	<-app.done
	defer close(app.done)

	fmt.Println("Shutting down")

	stopEventHandler <- true
	stopBroadcast <- true
	serverShutdown <- true

	if err := app.disconnect(); err != nil {
		log.Error(err)
	}
	time.Sleep(time.Second * shutdownDelay)
}

func (app *App) addEventListeners() {
	app.eventListener = &simconnect.EventListener{
		OnOpen:      app.OnOpen,
		OnQuit:      app.OnQuit,
		OnDataReady: app.OnDataReady,
		OnEventID:   app.OnEventID,
		OnException: app.OnException,
	}
}

func (app *App) initWebServer(address string, shutdown chan bool) {
	htmlHeaders := app.Headers(contentTypeHTML)
	textHeaders := app.Headers(contentTypeText)
	webServer := webserver.NewWebServer(address, shutdown)
	htmlDir := "assets/html"
	routes := []webserver.Route{
		{Pattern: "/", Handler: app.StaticContentHandler(htmlHeaders, "/", filepath.Join(htmlDir, "vfrmap.html"))},
		{Pattern: "/vfrmap", Handler: app.StaticContentHandler(htmlHeaders, "/vfrmap", filepath.Join(htmlDir, "vfrmap.html"))},
		{Pattern: "/mehmap", Handler: app.StaticContentHandler(htmlHeaders, "/mehmap", filepath.Join(htmlDir, "mehmap.html"))},
		{Pattern: "/setdata", Handler: app.StaticContentHandler(htmlHeaders, "/setdata", filepath.Join(htmlDir, "setdata.html"))},
		{Pattern: "/airports", Handler: app.StaticContentHandler(htmlHeaders, "/airports", filepath.Join(htmlDir, "airports.html"))},
		{Pattern: "/teleport", Handler: app.StaticContentHandler(htmlHeaders, "/teleport", filepath.Join(htmlDir, "teleporter.html"))},
		{Pattern: "/debug", Handler: app.GeneratedContentHandler(textHeaders, "/debug", app.DebugGenerator)},
		{Pattern: "/simvars", Handler: app.GeneratedContentHandler(textHeaders, "/simvars", app.SimvarsGenerator)},
		{Pattern: "/ws", Handler: app.socket.Serve},
	}

	fmt.Println("Starting web server")
	staticAssetsDir := "/assets/"
	webServer.Run(routes, staticAssetsDir)

	fmt.Printf("Web Server listening on %s\n\n", address)
	app.listNetworkInterfaces()
}

// https://golang-examples.tumblr.com/post/99458329439/get-local-ip-addresses
func (app *App) listNetworkInterfaces() {
	list, err := net.Interfaces()
	if err != nil {
		return
	}
	fmt.Println("Your network interfaces:")
	for i, iface := range list {
		str := fmt.Sprintf(" %d %s: ", i+1, iface.Name)
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for j, addr := range addrs {
			str += fmt.Sprintf("%v", addr)
			if j < len(addrs) {
				str += ", "
			}
		}
		fmt.Println(str)
	}
	fmt.Printf("\n")
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
			return fmt.Errorf("opening a connection to the simulator timed out")
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
			app.println("Received SIGTERM")
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
				app.println("Message", connID, msg)

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
		if app.airportsDB == nil {
			fmt.Println("Airports database not available")
			return
		}

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
		app.println("Added SimVar", defineID, name, unit, typ)
		request.Add(defineID, name, moniker)
	}, "data")
	app.requestManager.AddRequest(request)
	app.println("Added request", request)
}

func (app *App) handleSetDataMessage(msg *Message) {
	if !app.mate.IsConnected() {
		fmt.Println("Not connected to SimConnect. Ignoring SetDataMessage.")
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
		fmt.Println("Not connected to SimConnect. Ignoring TeleportMessage.")
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
				app.println("Removed SimVar", defineID)
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

// func (app *App) OnSimObjectData(data *simconnect.RecvSimObjectData) {
// 	// pass
// }

// func (app *App) OnSimObjectDataByType(data *simconnect.RecvSimObjectDataByType) {
// 	// pass
// }

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
		for j, simVar := range request.Vars {
			dump += fmt.Sprintf("    %02d: name: %s moniker: %s\n", j, simVar.Name, simVar.Moniker)
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

func (app *App) println(a ...interface{}) (n int, err error) {
	if app.verbose {
		return fmt.Println(a...)
	}
	return 0, nil
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
