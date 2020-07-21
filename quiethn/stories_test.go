package quiethn

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

const (
	fixturePath = "fixtures"
	itemPath    = fixturePath + "/items/%s"
)

func testServer() *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/topstories.json", func(w http.ResponseWriter, r *http.Request) {
		w.Write(renderJSON(fixturePath + "/topstories.json"))
	})
	mux.HandleFunc("/item/", func(w http.ResponseWriter, r *http.Request) {
		id := filepath.Base(r.URL.Path)
		path := fmt.Sprintf(itemPath, id)
		w.Write(renderJSON(path))
	})
	return httptest.NewServer(mux)
}

func renderJSON(path string) []byte {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	json, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	return json
}

func TestTopStories(t *testing.T) {
	oldNStories := nStories
	nStories = 3
	server := testServer()
	oldTopStoriesEndpoint := topStoriesEndpoint
	topStoriesEndpoint = server.URL + "/topstories.json"
	oldItemEndpoint := itemEndpoint
	itemEndpoint = server.URL + "/item/%d.json"

	defer func() {
		nStories = oldNStories
		topStoriesEndpoint = oldTopStoriesEndpoint
		itemEndpoint = oldItemEndpoint
		server.Close()
	}()

	stories, err := TopStories()
	if err != nil {
		t.Fatalf("received error: %v", err)
	}

	t.Run("returns the correct number of stories", func(t *testing.T) {
		if len(stories) != nStories {
			t.Errorf("got %d stories, want %d", len(stories), nStories)
		}
	})

	t.Run("returns 'story' items only", func(t *testing.T) {
		for _, story := range stories {
			if story.Type != "story" {
				t.Fatalf("item %d has type %q; only type 'story' is allowed", story.ID, story.Type)
			}
		}
	})

	t.Run("returns stories in the correct order", func(t *testing.T) {
		order := []int{1, 3, 4} // fixture 2 is of type 'comment' and should be skipped
		for i, story := range stories {
			fmt.Printf("%+v\n", story)
			if story.ID != order[i] {
				t.Fatalf("want story at position %d to have ID %d, got %d", i, order[i], story.ID)
			}
		}
	})
}
