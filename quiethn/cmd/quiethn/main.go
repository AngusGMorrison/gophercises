package main

import (
	"errors"
	"flag"
	"html/template"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/angusgmorrison/gophercises/quiethn"
)

type storyCache struct {
	stories []*quiethn.Item
	mu      sync.Mutex
}

var (
	tmpl                        *template.Template
	nStories                    = 30
	currentStories, nextStories *storyCache
)

func init() {
	// Initialize templates
	tmpl = template.Must(template.ParseGlob("templates/*"))
}

func main() {
	port := flag.String("port", "8080", "the port to listen on")
	flag.IntVar(&nStories, "nstories", 30, "the number of top stories to display")
	flag.Parse()

	initializeCaches()
	http.HandleFunc("/", showTopStories)
	log.Printf("Listening on port %s\n", *port)
	http.ListenAndServe("localhost:"+*port, nil)
}

func showTopStories(w http.ResponseWriter, r *http.Request) {
	stories, err := getCachedStories()
	if err != nil {
		log.Println(err)
		http.Error(w, "Unable to fetch top stories. Please try again later.",
			http.StatusInternalServerError)
		return
	}
	if err = tmpl.ExecuteTemplate(w, "top_stories", stories); err != nil {
		log.Println(err)
	}
}

func initializeCaches() {
	currentStories, nextStories = &storyCache{}, &storyCache{}
	cacheTopStories()
	rotateCache()
	go refresher()
}

func cacheTopStories() {
	log.Printf("caching nextStories...")
	stories, err := quiethn.TopStories(nStories)
	if err != nil {
		log.Printf("caching nextStories: %v\n", err)
		return
	}
	nextStories.mu.Lock()
	nextStories.stories = stories
	nextStories.mu.Unlock()
}

func getCachedStories() ([]*quiethn.Item, error) {
	currentStories.mu.Lock()
	defer currentStories.mu.Unlock()

	if currentStories.stories == nil {
		return nil, errors.New("Failed to load or rotate stories")
	}
	stories := make([]*quiethn.Item, len(currentStories.stories))
	copy(stories, currentStories.stories)
	return stories, nil
}

func rotateCache() {
	log.Printf("rotating cache...")
	nextStories.mu.Lock()
	currentStories.mu.Lock()

	currentStories.stories = nextStories.stories

	currentStories.mu.Unlock()
	nextStories.mu.Unlock()
}

const cacheDuration = 15

func refresher() {
	for {
		<-time.After((cacheDuration - 1) * time.Second)
		cacheTopStories()
		<-time.After(1 * time.Second)
		rotateCache()
	}
}
