package aeroports

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"strconv"
)

type Airport struct {
	id        int64
	Type      int
	ICAO      string
	Name      string
	Latitude  float64
	Longitude float64
	Elevation float64
	IATA      string
}

type Database struct {
	airports []*Airport
}

type Candidate struct {
	airport  *Airport
	distance float64
}

func NewDatabase() *Database {
	return &Database{
		airports: make([]*Airport, 0),
	}
}

const (
	degToRad    float64 = math.Pi / 180.0
	earthRadius float64 = 6371 * 1000.0
)

const (
	colID = iota
	colIdent
	colType
	colName
	colLatitude
	colLongitude
	colElevation
	colContinent
	colISOCountry
	colISORegion
	colMunicipality
	colScheduledService
	colGPSCode
	colIATACode
	colLocalCode
	colHomeLink
	colWikipediaLink
	colKeywords
)

const (
	AirportTypeUnknown  = 0x00
	AirportTypeClosed   = 0x01
	AirportTypeHeliport = 0x02
	AirportTypeSmall    = 0x04
	AirportTypeMedium   = 0x08
	AirportTypeLarge    = 0x10
	AirportTypeAll      = AirportTypeClosed | AirportTypeHeliport | AirportTypeSmall | AirportTypeMedium | AirportTypeLarge
	AirportTypeActive   = AirportTypeHeliport | AirportTypeSmall | AirportTypeMedium | AirportTypeLarge
	AirportTypeRunways  = AirportTypeSmall | AirportTypeMedium | AirportTypeLarge
)

func KilometersToMeters(km float64) float64 {
	return km * 1000.0
}

func MilestoMeters(mi float64) float64 {
	return mi * 1609.0
}

func NauticalMilesToMeters(nm float64) float64 {
	return nm * 1852.0
}

func (db *Database) ParseAirports(filename string, airportTypeFilter int, skipFirstLine bool) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	reader := csv.NewReader(f)

	line := -1
	for {
		row, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return err
		}

		line++

		if line == 0 && skipFirstLine {
			continue
		}

		typ := AirportTypeFromString(row[colType])
		if typ == AirportTypeUnknown {
			continue
		}
		if typ&airportTypeFilter == 0 {
			continue
		}

		id, err := parseInt(row[colID])
		if err != nil {
			fmt.Println(line, err)
			continue
		}

		latitude, err := parseFloat(row[colLatitude])
		if err != nil {
			fmt.Println(line, err)
			continue
		}

		longitude, err := parseFloat(row[colLongitude])
		if err != nil {
			fmt.Println(line, err)
			continue
		}

		elevation, err := parseFloat(row[colElevation])
		if err != nil {
			elevation = -1.0
		}

		icao := row[colIdent]
		name := row[colName]
		iata := row[colIATACode]

		airport := &Airport{
			id:        id,
			Type:      typ,
			ICAO:      icao,
			Name:      name,
			Latitude:  latitude,
			Longitude: longitude,
			Elevation: elevation,
			IATA:      iata,
		}
		db.airports = append(db.airports, airport)
	}
}

func (db *Database) List() {
	for i, airport := range db.airports {
		fmt.Println(i, *airport)
	}
}

func (db *Database) FindNearestAirport(latitude, longitude, radius float64, typeFilter int) *Airport {
	var nearestAirport *Airport = nil
	minDistance := math.MaxFloat64

	if radius < 0 {
		radius = math.MaxFloat64
	}

	for _, airport := range db.airports {
		if airport.Type&typeFilter == 0 {
			continue
		}
		d := distance(latitude, longitude, airport.Latitude, airport.Longitude)
		if d < minDistance && d <= radius {
			minDistance = d
			nearestAirport = airport
		}
	}
	return nearestAirport
}

func (db *Database) FindNearestAirports(latitude, longitude, radius float64, maxCount int, typeFilter int) []*Airport {
	candidates := make([]*Candidate, 0)

	if radius < 0 {
		radius = math.MaxFloat64
	}

	if maxCount < 0 {
		maxCount = math.MaxInt32
	}

	for _, airport := range db.airports {
		if airport.Type&typeFilter == 0 {
			continue
		}
		d := distance(latitude, longitude, airport.Latitude, airport.Longitude)
		if d <= radius {
			candidates = append(candidates, &Candidate{airport, d})
		}
	}

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].distance < candidates[j].distance
	})

	airports := make([]*Airport, 0, maxCount)
	count := min(len(candidates), maxCount)
	for i := 0; i < count; i++ {
		airports = append(airports, candidates[i].airport)
	}
	return airports
}

func AirportTypeFromString(typ string) int {
	switch typ {
	case "closed":
		return AirportTypeClosed
	case "small_airport":
		return AirportTypeSmall
	case "medium_airport":
		return AirportTypeMedium
	case "large_airport":
		return AirportTypeLarge
	case "heliport":
		return AirportTypeHeliport
	}
	return AirportTypeUnknown
}

func AirportTypeToString(typ int) string {
	switch typ {
	case AirportTypeClosed:
		return "closed"
	case AirportTypeHeliport:
		return "heliport"
	case AirportTypeSmall:
		return "small_airport"
	case AirportTypeMedium:
		return "medium_airport"
	case AirportTypeLarge:
		return "large_airport"
	}
	return "unknown"
}

func parseFloat(str string) (float64, error) {
	value, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0, err
	}
	return value, err
}

func parseInt(str string) (int64, error) {
	value, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, err
	}
	return value, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// see https://stackoverflow.com/questions/43167417/calculate-distance-between-two-points-in-leaflet
// returns the distance between to coordinates in meters
func distance(fromLatitude, fromLongitude, toLatitude, toLongitude float64) float64 {
	lat1 := fromLatitude * degToRad
	lon1 := fromLongitude * degToRad
	lat2 := toLatitude * degToRad
	lon2 := toLongitude * degToRad

	dtLat := lat2 - lat1
	dtLon := lon2 - lon1

	a := math.Pow(math.Sin(dtLat*0.5), 2) + math.Cos(lat1)*math.Cos(lat2)*math.Pow(math.Sin(dtLon*0.5), 2)
	c := 2 * math.Asin(math.Sqrt(a))
	return c * earthRadius
}
