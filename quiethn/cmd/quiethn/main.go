package main

import (
	"flag"
	"net/http"
)

func main() {
	port := flag.String("port", "8080", "the port to listen on")
	flag.Parse()

	http.HandleFunc("/", topStories)
	http.ListenAndServe(":"+*port, nil)
}

func topStories(w http.ResponseWriter, r *http.Request) {
	// Get top stories from the quiethn package
	// Render titles and URLS in a HTML template
}
