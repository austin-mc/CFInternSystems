package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// Base URL for the Cloudflare Radar endpoint
const radarEndpoint = "https://cfisysapi.developers.workers.dev"

// Used to hold the data from the Cloudflare Radar API endpoint
type RadarData struct {
	Timestamps []string
	Values     []string
}

// Used to return the stats calculated from the Cloudflare Radar API data
type RadarStats struct {
	Mean   string `json:"mean"`
	Median string `json:"median"`
	Min    string `json:"min"`
	Max    string `json:"max"`
}

func main() {
	// Setting up the mux router and HTTP handlers
	router := mux.NewRouter()

	router.HandleFunc("/", rootHandler)
	router.HandleFunc("/README.txt", readmeHandler)
	router.HandleFunc("/stats", statsHandler)

	http.ListenAndServe(":8080", router)
}

// Handles requests to "/" by returning a 404 Not Found
func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

// Handles requests to "/README.txt" by returning the README.txt file contents
func readmeHandler(w http.ResponseWriter, r *http.Request) {
	// Get bytes from README.txt
	data, err := getReadme()
	// Return 500 if there's an error - shouldn't ever happen
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Return a 200 response with the README text
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/text")
	w.Write(data)
}

// Gets the file contents of README.txt and returns them
func getReadme() ([]byte, error) {
	return ioutil.ReadFile("README.txt")
}

// Handles requests to "/stats" and responds with statistics from Cloudflare Radar
func statsHandler(w http.ResponseWriter, r *http.Request) {
	// Returns the user timestamp if provided, else returns empty string
	userTimestamp := r.URL.Query().Get("timestamp")
	// Returns the computed stats from the API data
	stats, err := prepareRadarRequest(userTimestamp)
	// Return 500 if there's an error
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json, err := json.Marshal(stats)
	// Return 500 if there's an error
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Respond with the stats
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

// Process the user timestamp if given and prepare the Radar API request URL
func prepareRadarRequest(queryTime string) (*RadarStats, error) {

	var timestamp int64
	var err error

	if len(queryTime) != 0 {
		// Try to convert query into int
		timestamp, err = strconv.ParseInt(queryTime, 10, 64)
		if err != nil {
			timestamp = time.Now().Unix()
		} else if timestamp > time.Now().Unix() {
			newError := errors.New("timestamp is in the future")
			return nil, newError
		}

	} else {
		timestamp = time.Now().Unix()
	}

	// Formatting the final request URL
	requestUrl := fmt.Sprintf("%s/stats?timestamp=%d", radarEndpoint, timestamp)
	// Get the Radar data
	return getRadarData(requestUrl)
}

// Gets the data from the Cloudflare Radar API
func getRadarData(url string) (*RadarStats, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var data RadarData
	json.Unmarshal(body, &data)
	return calculateStats(data)
}

// Uses the data from the Cloudflare Radar API to calculate the mean, median, min, and max
func calculateStats(data RadarData) (*RadarStats, error) {

	// Edge case where values struct is empty
	if len(data.Values) == 0 {
		return &RadarStats{Mean: "0.000", Median: "0.000", Min: "0.000", Max: "0.000"}, nil
	}

	// Making an int slice to convert strings
	values := make([]float64, len(data.Values))

	// Converting strings to int and calculating total
	total := 0.0
	for i, s := range data.Values {
		val, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return nil, err
		}
		values[i] = val
		total += val
	}

	// Sorting to find median
	sort.Float64s(values)

	// Calculating stats
	mean := total / float64(len(values))
	min := values[0]
	max := values[len(values)-1]
	length := len(values)

	var median float64

	if length%2 == 0 {
		// Even length array, using the average of the 2 median numbers
		median = (values[(length-1)/2] + values[(length+1)/2]) / 2
	} else {
		// Odd length array, just get the middle number
		median = values[length/2]
	}

	// Converting values into formatted strings and placing in struct for response
	stats := RadarStats{
		Mean:   fmt.Sprintf("%.3f", mean),
		Median: fmt.Sprintf("%.3f", median),
		Min:    fmt.Sprintf("%.3f", min),
		Max:    fmt.Sprintf("%.3f", max),
	}

	return &stats, nil
}
