package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/angusgmorrison/gophercises/url_shortener/urlshort"
)

var defaultRedirects = map[string]string{
	"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
	"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
}

func main() {
	format := flag.String(
		"format",
		"",
		"format of the optional file arg mapping short URLs to destinations. Supported formats: json, yaml",
	)
	flag.Parse()

	// Build the default MapHandler using the mux as the fallback
	mux := defaultMux()
	handler := urlshort.MapHandler(defaultRedirects, mux)

	// Read URL mapping from an optional JSON or YAML file
	var fileRedirects []byte
	if *format != "" {
		bytes, err := readFileData()
		if err != nil {
			exit(err.Error())
		}
		fileRedirects = bytes
	}

	switch *format {
	case "yaml":
		YAMLHandler, err := urlshort.YAMLHandler(fileRedirects, handler)
		if err != nil {
			exit(fmt.Sprintf("%s: urlshort.YAMLHandler: %v", os.Args[0], err))
		}
		handler = YAMLHandler
	case "json":
		JSONHandler, err := urlshort.JSONHandler(fileRedirects, handler)
		if err != nil {
			exit(fmt.Sprintf("%s: urlshort.JSONHandler: %v", os.Args[0], err))
		}
		handler = JSONHandler
	}

	log.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", handler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

func readFileData() ([]byte, error) {
	var file *os.File
	args := flag.Args()
	if len(args) == 0 {
		return nil, fmt.Errorf("usage: %s -format=[json|yaml] <file path>", os.Args[1])
	}
	file, err := os.Open(flag.Args()[0])
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(file)
}

func exit(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}
