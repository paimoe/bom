package main

var Locations map[string]string
var Stations map[string]map[string]int
var PostCodes map[int]string

// location shit cause idk how to import
func SearchLocation(label string) (map[string]int, bool) {
	loc, ok := Stations[label]
	return loc, ok
}

func map_location(x, y, station_num int) map[string]int {
	innerMap := make(map[string]int)
	innerMap["x"] = x
	innerMap["y"] = y
	innerMap["station_num"] = station_num
	return innerMap
}

func InitLocations() {
	Locations = make(map[string]string)
	Locations["Perth"] = "bwa_pt053"
	Locations["Sydney"] = "bnsw_pt131"
	Locations["Darwin"] = "bnt_pt001"
	Locations["Melbourne"] = "bvic_pt042"
	Locations["Brisbane"] = "bqld_pt001"
	Locations["Adelaide"] = "bsa_pt001"
	Locations["Hobart"] = "btas_pt021"
	Locations["Canberra"] = "bnsw_pt027"

	// These are hardcoded from https://api.bom.gov.au/apikey/v1/locations/places/details/place/btas_pt021?filter=nearby_type%3Abom_stn&radius=100000
	// where btas_pt021 is from above
	// This file will have a place.locationHierarchy.nearest.gridcells.forecast with an x and y, ie Hobart is x=593, y=42
	// This is passed into https://api.bom.gov.au/apikey/v1/forecasts/daily/{ x }/{ y }?timezone=Australia/{ tz }
	// which has fcst.daily.[0].atm.surf_air.temp_max_cel and temp_min_cel
	Stations = make(map[string]map[string]int)
	Stations["Perth"] = map_location(65, 262, 9225)
}
