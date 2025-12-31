package main

import (
	//"os"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"time"

	"github.com/gorilla/mux"

	// "io/ioutil"
	"github.com/allegro/bigcache/v3"
)

var Locations map[string]string
var Stations map[string]map[string]int
var PostCodes map[int]string

var cache *bigcache.BigCache

const Port = 8190

// location shit cause idk how to import
func searchLocation(label string) string {
	return Locations[label]
}

func map_location(x, y int) map[string]int {
	innerMap := make(map[string]int)
	innerMap["x"] = x
	innerMap["y"] = y
	return innerMap
}

func init_locations() {
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
	Stations["Perth"] = map_location(65, 262)
}

func getAPIResponse(url string) (*Response, error) {
	cache_key := url
	// Check cache first
	entry, err := cache.Get(cache_key)
	if err == nil {
		// Cache hit
		fmt.Println("Cache hit for", cache_key)
		var cachedResponse Response
		err = json.Unmarshal(entry, &cachedResponse)
		if err == nil {
			return &cachedResponse, nil
		}
	}

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:132.0) Gecko/20100101 Firefox/132.0")
	// req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	// req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("Referer", "https://www.bom.gov.au/")
	req.Header.Set("Origin", "https://www.bom.gov.au/")
	// req.Header.Set("DNT", "1")
	// req.Header.Set("Sec-GPC", "1")
	// req.Header.Set("Connection", "keep-alive")
	// req.Header.Set("Sec-Fetch-Dest", "empty")
	// req.Header.Set("Sec-Fetch-Mode", "cors")
	// req.Header.Set("Sec-Fetch-Site", "same-site")
	// req.Header.Set("Priority", "u=4")
	// req.Header.Set("TE", "trailers")

	fmt.Println("Making request to:", url)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make get request: %v", err)
	}
	defer resp.Body.Close()

	// bodyBytes, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to read body: %v", err)
	// }

	// // Print the raw string
	// bodyString := string(bodyBytes)
	// fmt.Println(bodyString)

	var result Response
	decoder := json.NewDecoder(resp.Body)
	decoder.DisallowUnknownFields() // optional but useful
	err = decoder.Decode(&result)

	if err != nil {
		return nil, fmt.Errorf("failed to decode json: %v", err)
	}

	// Store in cache
	responseBytes, err := json.Marshal(result)
	if err == nil {
		cache.Set(cache_key, responseBytes)
	}

	return &result, nil
}

func main() {
	// we want to listen to requests and respond with either cached
	// bom data, or update the cache and return.
	init_locations()

	// Cache
	cache, _ = bigcache.New(context.Background(), bigcache.DefaultConfig(60*time.Minute))

	r := mux.NewRouter()
	r.HandleFunc("/{location}", LocationHandler)
	r.HandleFunc("/", HomeHandler)

	// Run server
	http.Handle("/", r)

	fmt.Printf("Starting on port %d", Port)
	http.ListenAndServe(fmt.Sprintf(":%d", Port), nil)
}

func LocationHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	loc := vars["location"]

	match, _ := regexp.MatchString("^[0-9]{4}$", loc)

	if match {
		// is by postcode. Pick a useful default like 6164 = Cockburn Central
	} else {
		// Match by location. Check in Locations list
		location_code := searchLocation(loc)
		if location_code == "" {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "404 Location not found: %s", loc)
			// TODO try these similar names
			return
		}

		// Get from api (no cache yet)
		// Step 1. Get nearest weather station
		station_coords := Stations[loc]
		url := fmt.Sprintf("https://api.bom.gov.au/apikey/v1/forecasts/daily/%d/%d?timezone=Australia/Perth", station_coords["x"], station_coords["y"])

		response, err := getAPIResponse(url)
		if err != nil {
			fmt.Println("Error getting API response:", err)
		}

		fmt.Fprintf(w, "API Response: %+v", response)
	}
}

type TemplateData struct {
	Capitals []string
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	// add list of locations to template NOT USED
	location_keys := make([]string, 0, len(Locations))
	for k := range Locations {
		location_keys = append(location_keys, k)
	}

	data := TemplateData{Capitals: location_keys}

	t, _ := template.ParseFiles("static/index.html")
	t.Execute(w, data)
}
