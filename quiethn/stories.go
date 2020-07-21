package quiethn

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

var (
	nStories           = 30
	baseEndpoint       = "https://hacker-news.firebaseio.com/v0"
	topStoriesEndpoint = baseEndpoint + "/topstories.json"
	itemEndpoint       = baseEndpoint + "/item/%d.json"
)

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
func TopStories() ([]*Item, error) {
	itemIDs, err := getTopItemIDs()
	if err != nil {
		return nil, fmt.Errorf("TopStories(): %v", err)
	}

	stories, err := getTopStories(itemIDs)
	if err != nil {
		return nil, fmt.Errorf("TopStories(): %v", err)
	}

	return stories, nil
}

// getTopItemIDs returns the ids for the top 500 stories currently on
// Hacker News.
func getTopItemIDs() ([]int, error) {
	resp, err := http.Get(topStoriesEndpoint)
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

// getTopStories fetches the data for each item ID, filtering out
// non-story items until maxStories stories have been retrieved.
func getTopStories(IDs []int) ([]*Item, error) {
	stories := make([]*Item, 0)
	for _, id := range IDs {
		if len(stories) == nStories {
			break
		}

		story, err := getStory(id)
		if err != nil {
			if err != errItemType {
				fmt.Fprintln(os.Stderr, err)
			}
			continue
		}
		stories = append(stories, story)
	}

	return stories, nil
}

var errItemType = errors.New("item is not a story")

func getStory(id int) (*Item, error) {
	itm, err := getItem(id)
	if err != nil {
		return nil, err
	}
	if itm.Type != "story" {
		return nil, errItemType
	}
	return itm, nil
}

func getItem(id int) (*Item, error) {
	endpoint := fmt.Sprintf(itemEndpoint, id)
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
