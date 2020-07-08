package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/angusgmorrison/gophercises/url_shortener/from_db/urlshort"
)

var seed bool

func main() {
	flag.BoolVar(
		&seed,
		"seed",
		false,
		"seed the DB with data from the YAML file specified as an argument",
	)
	flag.Parse()

	// Initialize the DB
	store, err := initializeStore()
	if err != nil {
		exit(err.Error())
	}
	defer store.DB.Close()

	// Build the default MapHandler using the default mux as the fallback.
	mux := defaultMux()
	handler := urlshort.Handler(store, mux)

	log.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", handler)
}

func initializeStore() (*urlshort.RedirectStore, error) {
	store, err := urlshort.OpenRedirectStore("redirects.db")
	if err != nil {
		return nil, fmt.Errorf("opening DB: %v", err)
	}

	// Seed the DB
	if seed {
		args := flag.Args()
		if len(args) == 0 {
			return nil, fmt.Errorf("usage: %s -format=[json|yaml] <file path>", os.Args[0])
		}
		if err := store.Seed(args[0]); err != nil {
			return nil, fmt.Errorf("seeding DB: %v", err)
		}
	}
	return store, err
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

func exit(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}
