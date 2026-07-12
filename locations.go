package main

import (
	"strings"
)

type Location struct {
	Label      string
	X          int
	Y          int
	StationNum int
}

var Locations map[string]string
var Stations map[string]Location
var PostCodes map[int]string

func SearchLocation(label string) (Location, bool) {
	loc, ok := Stations[strings.ToLower(label)]
	return loc, ok
}

func InitLocations() {
	Locations = make(map[string]string)
	Locations["perth"] = "bwa_pt053"
	Locations["sydney"] = "bnsw_pt131"
	Locations["darwin"] = "bnt_pt001"
	Locations["melbourne"] = "bvic_pt042"
	Locations["brisbane"] = "bqld_pt001"
	Locations["adelaide"] = "bsa_pt001"
	Locations["hobart"] = "btas_pt021"
	Locations["canberra"] = "bnsw_pt027"

	// These are hardcoded from https://api.bom.gov.au/apikey/v1/locations/places/details/place/btas_pt021?filter=nearby_type%3Abom_stn&radius=100000
	// where btas_pt021 is from above
	// This file will have a place.locationHierarchy.nearest.gridcells.forecast with an x and y, ie Hobart is x=593, y=42
	// This is passed into https://api.bom.gov.au/apikey/v1/forecasts/daily/{ x }/{ y }?timezone=Australia/{ tz }
	// which has fcst.daily.[0].atm.surf_air.temp_max_cel and temp_min_cel
	Stations = make(map[string]Location)
	Stations["perth"] = Location{Label: "Perth", X: 65, Y: 262, StationNum: 9225}
	Stations["sydney"] = Location{Label: "Sydney", X: 658, Y: 223, StationNum: 66214}
	Stations["darwin"] = Location{Label: "Darwin", X: 316, Y: 652, StationNum: 14318}
	Stations["melbourne"] = Location{Label: "Melbourne", X: 553, Y: 143, StationNum: 86338}
	Stations["brisbane"] = Location{Label: "Brisbane", X: 688, Y: 349, StationNum: 40211}
	Stations["adelaide"] = Location{Label: "Adelaide", X: 445, Y: 201, StationNum: 23154}
	Stations["hobart"] = Location{Label: "Hobart", X: 591, Y: 42, StationNum: 94087}
	Stations["canberra"] = Location{Label: "Canberra", X: 624, Y: 194, StationNum: 70351}
}
