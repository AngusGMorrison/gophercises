package quiethn

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	baseEndpoint       = "https://hacker-news.firebaseio.com/v0"
	topStoriesEndpoint = "/topstories.json"
	itemEndpoint       = "/item/"
)

const maxTopStories = 500

// Item represents a HackerNews item without its text content and
// fields related to comments and polls.
type Item struct {
	ID        int    `json:"id"`
	Deleted   bool   `json:"deleted"`
	Type      string `json:"type"`
	By        string `json:"by"`
	CreatedAt int    `json:"time"`
	Dead      bool   `json:"dead"`
	URL       string `json:"url"`
	Score     int    `json:"score"`
	Title     string `json:"title"`
}

// TopStories fetches the top maxStories stories from HackerNews,
// ignoring comments, users and job postings.
func TopStories(nStories int) ([]*Item, error) {
	itemIDs, err := TopItemIDs()
	if err != nil {
		return nil, fmt.Errorf("TopStories(): %v", err)
	}

	nStories = min(nStories, maxTopStories) // prevent a request exceeding the max possible stories
	stories, err := getTopStories(itemIDs, nStories)
	if err != nil {
		return nil, fmt.Errorf("TopStories(): %v", err)
	}

	return stories, nil
}

// TopItemIDs returns the ids for the top 500 stories currently on
// Hacker News.
func TopItemIDs() ([]int, error) {
	resp, err := http.Get(fmt.Sprintf("%s%s", baseEndpoint, topStoriesEndpoint))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var itemIDs []int
	err = json.NewDecoder(resp.Body).Decode(&itemIDs)
	if err != nil {
		return nil, fmt.Errorf("decoding response: %v", err)
	}
	return itemIDs, nil
}

type result struct {
	Story *Item
	Err   error
}

// getTopStories fetches the data for each item ID, filtering out
// non-story items until maxStories stories have been retrieved. In
// the event that more stories are requested than there are available,
// the meaximum number of available stories is returned.
func getTopStories(IDs []int, max int) ([]*Item, error) {
	var pending, skipped int
	var seenAll bool
	stories := make(map[int]*Item)
	sema := make(chan struct{}, 20) // counting semaphore
	results := make(chan *result)

	for len(stories) < max && !seenAll { // until sufficient stories are found...
		// ...start enough goroutines to populate the remaining items
		// on the assumption that all will be stories.
		found := len(stories)
		for i := 0; i < max-found; i++ {
			nextID := found + skipped + i
			if nextID == len(IDs) {
				seenAll = true
				break
			}
			pending++
			go getStoryConcurrent(IDs[nextID], sema, results)
		}

		// While goroutines are active, receive from them and populate
		// the stories map.
		for ; pending > 0; pending-- {
			result := <-results
			if result.Err != nil {
				skipped++
			} else {
				stories[result.Story.ID] = result.Story
			}
		}
	}

	IDsProcessed := len(stories) + skipped
	return orderStories(IDs, stories, IDsProcessed), nil
}

var errItemType = errors.New("item is not a story")

func getStoryConcurrent(id int, sema chan struct{}, results chan<- *result) {
	sema <- struct{}{}        // acquire a token
	defer func() { <-sema }() // release token when done

	story, err := GetStory(id)
	if err != nil {
		results <- &result{Err: err}
		return
	}
	results <- &result{Story: story}
}

// GetStory fetches an item from Hacker News given its ID and returns
// it if it has Type == "story", or an error otherwise.
func GetStory(id int) (*Item, error) {
	itm, err := GetItem(id)
	if err != nil {
		return nil, err
	}
	if itm.Type != "story" {
		return nil, errItemType
	}
	return itm, nil
}

// GetItem fetches and returns a single item from Hacker News given
// its ID.
func GetItem(id int) (*Item, error) {
	endpoint := fmt.Sprintf("%s%s%d.json", baseEndpoint, itemEndpoint, id)
	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("fetching item %d: %v", id, err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading item %d: %v", id, err)
	}

	var itm Item
	err = json.Unmarshal(body, &itm)
	if err != nil {
		return nil, fmt.Errorf("decoding item %d: %v", id, err)
	}
	return &itm, nil
}

func orderStories(IDs []int, stories map[int]*Item, IDsProcessed int) []*Item {
	ordered := make([]*Item, len(stories))
	for i, j := 0, 0; i < IDsProcessed; i++ {
		storyID := IDs[i]
		if story, ok := stories[storyID]; ok {
			ordered[j] = story
			j++
		}
	}
	return ordered
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
