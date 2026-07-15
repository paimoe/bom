package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"maps"
	"slices"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/gorilla/mux"
)

// Helper for env var lookup, ensures compatibility with testing
func lookupEnv(key string) (string, bool) {
	return os.LookupEnv(key)
}

var cache *bigcache.BigCache

// Port is now configured using the PORT environment variable, falls back to 8190 if not set
func getPort() int {
	port := 8190
	if val, found := lookupEnv("PORT"); found {
		if parsed, err := strconv.Atoi(val); err == nil {
			port = parsed
		}
	}
	return port
}

func getAPIResponse(url string, dest interface{}) (error, bool) {
	cache_key := url
	// Check cache first
	entry, err := cache.Get(cache_key)
	if err == nil {
		// Cache hit
		fmt.Println("Cache hit for", cache_key)
		err = json.Unmarshal(entry, dest)
		if err == nil {
			return nil, true
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
		return fmt.Errorf("failed to make get request: %v", err), false
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	decoder.DisallowUnknownFields() // optional but useful
	err = decoder.Decode(dest)
	if err != nil {
		return fmt.Errorf("failed to decode json: %v", err), false
	}

	// Store in cache
	responseBytes, err := json.Marshal(dest)
	if err == nil {
		cache.Set(cache_key, responseBytes)
	}

	return nil, false
}

type HomeTemplateData struct {
	Capitals []string
}
type PlaceTemplateData struct {
	Location  string
	TempNow   float64
	FeelsLike float64
	TodayHigh float64
	TodayLow  float64
	IsCached  bool
}

func main() {
	// we want to listen to requests and respond with either cached
	// bom data, or update the cache and return.
	InitLocations()

	// Cache
	cache, _ = bigcache.New(context.Background(), bigcache.DefaultConfig(20*time.Minute))

	r := mux.NewRouter()
	r.HandleFunc("/{location}", LocationHandler)
	r.HandleFunc("/", HomeHandler)

	// Handle static css
	fs := http.FileServer(http.Dir("./static/assets"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	// Run server
	http.Handle("/", r)

	port := getPort()
	fmt.Printf("Starting on port %d\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func LocationHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	loc := vars["location"]

	match, _ := regexp.MatchString("^[0-9]{4}$", loc)

	if match {
		// is by postcode. Pick a useful default like 6164 = Cockburn Central
	} else {
		// Match by location. Check in Locations list
		location, ok := SearchLocation(loc)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "404 Location not found: %s", loc) // Maybe log it a bit
			// TODO try these similar names
			return
		}

		// Get from api (no cache yet)
		// Step 1. Get nearest weather station
		url_fcast := fmt.Sprintf("https://api.bom.gov.au/apikey/v1/forecasts/daily/%d/%d?timezone=Australia/Perth", location.X, location.Y)
		url_now := fmt.Sprintf("https://api.bom.gov.au/apikey/v1/observations/latest/%d/atm/surf_air?include_qc_results=false", location.StationNum)

		var response_now ResponseNow
		err, cached_now := getAPIResponse(url_now, &response_now)
		if err != nil {
			fmt.Println("Error getting API Now response:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var response_fcast ResponseFcast
		err, cached_fcast := getAPIResponse(url_fcast, &response_fcast)
		if err != nil {
			fmt.Println("Error getting API Forecast response:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// fmt.Fprintf(w, "Temperature: %f", *response_fcast.Fcst.Daily[1].Atm.SurfAir.TempMaxCel)
		// fmt.Fprintf(w, "API Response: %+v", response_fcast)
		// fmt.Fprintf(w, "API Response: %+v", response_now)

		tpl_data := PlaceTemplateData{
			Location:  location.Label,
			TempNow:   response_now.Obs.Temp.DryBulb1MinCel,
			FeelsLike: response_now.Obs.Temp.Apparent1MinCel,
			TodayHigh: response_now.Obs.Temp.DryBulbMaxCel,
			TodayLow:  response_now.Obs.Temp.DryBulbMinCel,
			IsCached:  cached_now && cached_fcast,
		}
		t, err := template.ParseFiles("static/_base.html", "static/place.html")
		if err != nil {
			http.Error(w, "Internal Server Error: loading template", http.StatusInternalServerError)
			fmt.Println("Template error:", err)
			return
		}
		t.ExecuteTemplate(w, "base", tpl_data)
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	// add list of locations to template NOT USED
	location_keys := slices.Collect(maps.Keys(Locations))

	data := HomeTemplateData{Capitals: location_keys}

	t, _ := template.ParseFiles("static/_base.html", "static/index.html")
	t.ExecuteTemplate(w, "base", data)
}
