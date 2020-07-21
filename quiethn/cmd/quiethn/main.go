package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"

	"github.com/angusgmorrison/gophercises/quiethn"
)

var (
	tmpl *template.Template
)

func init() {
	// Initialize templates
	tmpl = template.Must(template.ParseGlob("templates/*"))
}

func main() {
	port := flag.String("port", "8080", "the port to listen on")
	flag.Parse()

	http.HandleFunc("/", topStories)
	log.Printf("Listening on port %s\n", *port)
	http.ListenAndServe("localhost:"+*port, nil)
}

func topStories(w http.ResponseWriter, r *http.Request) {
	stories, err := quiethn.TopStories(30)
	if err != nil {
		http.Error(w, "Unable to fetch top stories. Please try again later.",
			http.StatusInternalServerError)
		return
	}
	if err = tmpl.ExecuteTemplate(w, "top_stories", stories); err != nil {
		log.Println(err)
	}
}
