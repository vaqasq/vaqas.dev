package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func fetchImage(apiKey string) string {

	// define struct to get the field you want
	type Response struct {
		ImageURL string `json:"url"`
	}

	//api
	reponse, err := http.Get("https://api.nasa.gov/planetary/apod?api_key=" + apiKey)

	if err != nil {
		log.Fatal(err)
	}

	defer reponse.Body.Close()

	var resp Response

	if err := json.NewDecoder(reponse.Body).Decode(&resp); err != nil {
		log.Fatal(err)
	}

	// resp.ImageURL is the image URL
	return resp.ImageURL

}

func handler(w http.ResponseWriter, r *http.Request) {

	apiKey := os.Getenv("NASA_API_KEY")
	imageURL := fetchImage(apiKey)

	tmpl, err := template.ParseFiles("static-files/index.html")

	if err != nil {
		log.Fatal(err)
	}

	tmpl.Execute(w, imageURL)

}

func main() {

	godotenv.Load()

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
