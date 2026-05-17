package main

import (
	"log"
	"net/http"
)

/*
func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, world!")
}
*/

func main() {
	http.Handle("/", http.FileServer(http.Dir("./static-files")))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
