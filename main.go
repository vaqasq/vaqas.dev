package main

import (
	"encoding/json"
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
	reponse, err := http.Get("https://api.nasa.gov/planetary/apod?api_key=" + apiKey)

	if err != nil {
		return nil, err
	}

	defer reponse.Body.Close()

	var resp Response

	if err := json.NewDecoder(reponse.Body).Decode(&resp); err != nil {
		return nil, err
	}

	// resp.ImageURL is the image URL
	return []string{resp.ImageURL, resp.MediaType}, nil

}

func handler(w http.ResponseWriter, r *http.Request) {

	apiKey := os.Getenv("NASA_API_KEY")

	// Current UTC date to be used as a key for the cache
	dateOnly := time.Now().UTC().Format("2006-01-02")
	raw, found := C.Get(dateOnly)

	var fetchSlice []string

	if found {

		var ok bool
		fetchSlice, ok = raw.([]string)

		if ok == false {
			var err error
			fetchSlice, err = fetchImage(apiKey)

			if err != nil {
				http.Error(w, "Failed to fetch image", http.StatusInternalServerError)
				return
			}

			C.Set(dateOnly, fetchSlice, cache.DefaultExpiration)

		}

	} else {

		var err error
		fetchSlice, err = fetchImage(apiKey)

		if err != nil {
			http.Error(w, "Failed to fetch image", http.StatusInternalServerError)
			return
		}

		C.Set(dateOnly, fetchSlice, cache.DefaultExpiration)
	}

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
