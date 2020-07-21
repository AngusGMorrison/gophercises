package quiethn

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

const (
	fixturePath = "fixtures"
	itemPath    = fixturePath + "/items/%s"
)

func testServer() *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/topstories.json", func(w http.ResponseWriter, r *http.Request) {
		data, _ := renderJSON(fixturePath + "/topstories.json")
		w.Write(data)
	})
	mux.HandleFunc("/item/", func(w http.ResponseWriter, r *http.Request) {
		id := filepath.Base(r.URL.Path)
		path := fmt.Sprintf(itemPath, id)
		data, err := renderJSON(path)
		if err != nil {
			http.Error(w, "Page not found", http.StatusNotFound)
			return
		}
		w.Write(data)
	})
	return httptest.NewServer(mux)
}

func renderJSON(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err // report as 404
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err) // should not happen in test environment
	}
	return data, nil
}

func TestTopStories(t *testing.T) {
	server := testServer()
	oldTopStoriesEndpoint := topStoriesEndpoint
	topStoriesEndpoint = server.URL + "/topstories.json"
	oldItemEndpoint := itemEndpoint
	itemEndpoint = server.URL + "/item/%d.json"

	defer func() {
		topStoriesEndpoint = oldTopStoriesEndpoint
		itemEndpoint = oldItemEndpoint
		server.Close()
	}()

	maxStories := 3
	stories, err := TopStories(maxStories)
	if err != nil {
		t.Fatalf("received error: %v", err)
	}

	t.Run("returns the correct number of stories", func(t *testing.T) {
		if len(stories) != maxStories {
			t.Errorf("got %d stories, want %d", len(stories), maxStories)
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
			if story.ID != order[i] {
				t.Fatalf("want story at position %d to have ID %d, got %d", i, order[i], story.ID)
			}
		}
	})
}

func TestTopItemIDs(t *testing.T) {
	server := testServer()
	oldEndpoint := topStoriesEndpoint
	topStoriesEndpoint = server.URL + "/topstories.json"

	defer func() {
		server.Close()
		topStoriesEndpoint = oldEndpoint
	}()

	t.Run("returns the correct item IDs in order", func(t *testing.T) {
		IDs, err := TopItemIDs()
		if err != nil {
			t.Fatalf("received error: %v", err)
		}
		for i, gotID := range IDs {
			wantID := i + 1
			if gotID != wantID {
				t.Fatalf("got ID %d, want %d", gotID, wantID)
			}
		}
	})
}

func TestGetStory(t *testing.T) {
	server := testServer()
	oldEndpoint := itemEndpoint
	itemEndpoint = server.URL + "/item/%d.json"

	defer func() {
		server.Close()
		itemEndpoint = oldEndpoint
	}()

	t.Run("returns an item when the ID corresponds to a story", func(t *testing.T) {
		wantID := 1 // fixture 1 is a story
		itm, err := GetStory(wantID)
		if err != nil {
			t.Fatalf("GetStory(%d): recevied error: %v", wantID, err)
		}
		if itm.ID != wantID {
			t.Errorf("GetStory(%d): got story with ID %d, want %d", wantID, itm.ID, wantID)
		}
		if itm.Type != "story" {
			t.Errorf("GetStory(%d): got Item with type %q, want \"story\"", wantID, itm.Type)
		}
	})

	t.Run("returns errItemType when the ID does not correspond to a story", func(t *testing.T) {
		wantID := 2 // fixture 2 is a comment
		_, err := GetStory(wantID)
		if err == nil {
			t.Fatalf("GetStory(%d): failed to return error with non-story ID", wantID)
		}
		if err != errItemType {
			t.Fatalf("GetStory(%d): got error %q, want errItemType", wantID, err)
		}
	})
}

func TestGetItem(t *testing.T) {
	server := testServer()
	oldEndpoint := itemEndpoint
	itemEndpoint = server.URL + "/item/%d.json"

	defer func() {
		server.Close()
		itemEndpoint = oldEndpoint
	}()

	t.Run("returns the correct Item when the ID is valid", func(t *testing.T) {
		wantID := 1
		wantItm := Item{
			ID:        wantID,
			Deleted:   false,
			By:        "scastiel",
			Score:     292,
			CreatedAt: 1595245898,
			Title:     "Show HN: 3D Book Image CSS Generator",
			Type:      "story",
			URL:       "https://3d-book-css.netlify.app/",
		}
		gotItm, err := GetItem(wantID)
		if err != nil {
			t.Fatalf("GetItem(%d): receieved error: %v", wantID, err)
		}
		if !reflect.DeepEqual(wantItm, *gotItm) {
			t.Fatalf("GetItem(%d): got %+v, want %+v", wantID, gotItm, wantItm)
		}
	})

	t.Run("returns an error when the ID is invalid", func(t *testing.T) {
		wantID := 100
		if _, err := GetItem(wantID); err == nil {
			t.Fatalf("GetItem(%d): failed to return error with invalid ID", wantID)
		}
	})

}
