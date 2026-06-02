package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/patrickmn/go-cache"
)

// Cache as a global variable. Don't need a closure for the handler now
var C = cache.New(26*time.Hour, 6*time.Hour)

func fetchImage(apiKey string) ([]string, error) {

	// define struct to get the field you want
	type Response struct {
		ImageURL  string `json:"url"`
		MediaType string `json:"media_type"`
	}

	//api
	response, err := http.Get("https://api.nasa.gov/planetary/apod?api_key=" + apiKey)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("NASA API failed with status: %d", response.StatusCode)
	}

	var resp Response

	if err := json.NewDecoder(response.Body).Decode(&resp); err != nil {
		return nil, err
	}

	// resp.ImageURL is the image URL
	return []string{resp.ImageURL, resp.MediaType}, nil

}

func handler(w http.ResponseWriter, r *http.Request) {

	apiKey := os.Getenv("NASA_API_KEY")
	dateOnly := time.Now().UTC().Format("2006-01-02")

	raw, found := C.Get(dateOnly)
	var fetchSlice []string

	// Creates a flag to track if we need to call the API
	needsFetch := false

	if found {
		var ok bool
		fetchSlice, ok = raw.([]string)

		if !ok {
			needsFetch = true
		}
	} else {
		needsFetch = true
	}

	if needsFetch {
		var err error
		fetchSlice, err = fetchImage(apiKey)

		if err != nil {
			http.Error(w, "Failed to fetch image", http.StatusInternalServerError)
			return
		}

		// Calculate time until midnight, aligns with APOD
		now := time.Now().UTC()
		nextMidnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, time.UTC)
		timeUntilMidnight := nextMidnight.Sub(now)

		C.Set(dateOnly, fetchSlice, timeUntilMidnight)
	}

	// Render the template
	data := struct {
		ImageURL  string
		MediaType string
	}{
		ImageURL:  fetchSlice[0],
		MediaType: fetchSlice[1],
	}

	tmpl, err := template.ParseFiles("static-files/index.html")

	if err != nil {
		http.Error(w, "Failed to parse static files ", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Failed to execute templates", http.StatusInternalServerError)
		return
	}
}

func main() {

	godotenv.Load()

	http.HandleFunc("/", handler)
	http.Handle("/static-files/", http.StripPrefix("/static-files/", http.FileServer(http.Dir("./static-files"))))
	log.Fatal(http.ListenAndServe(":8080", nil))

}
