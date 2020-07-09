package cyoa

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

func TestJSONStory(t *testing.T) {
	wantChapter := Chapter{
		Title:      "C1",
		Paragraphs: []string{"This is Chapter 1."},
		Options: []Option{
			{"Gather your party", "/venture-forth/"},
		},
	}
	chapter2Title := "C2"
	fixture := `{
		"%[1]s": {
			"title": "%[1]s",
			"story": [
				"This is Chapter 1."
			],
			"options": [
				{
					"text": "%[2]s",
					"chapter": "%[3]s"
				}
			]
		},
		"%[4]s": {
			"title": "%[4]s",
			"story": [
				"This is Chapter 2."
			],
			"options": [
				{
					"text": "Option text",
					"chapter": "/option/"
				}
			]
		}
	}`

	// Parse fixture
	jsonString := fmt.Sprintf(fixture, wantChapter.Title, wantChapter.Options[0].Text,
		wantChapter.Options[0].Chapter, chapter2Title)
	r := bytes.NewReader([]byte(jsonString))
	story, err := JSONStory(r)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// Test existence of chapters
	for _, chapter := range []string{wantChapter.Title, chapter2Title} {
		if _, ok := story[chapter]; !ok {
			t.Fatalf(err.Error())
		}
	}

	// Test correctness of first chapter
	gotChapter := story[wantChapter.Title]
	if !reflect.DeepEqual(gotChapter, wantChapter) {
		t.Fatalf("JSONStory returned chapter 1 %+v, want %+v", gotChapter, wantChapter)
	}
}
