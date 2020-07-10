package cyoa

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"
)

const storyFixture = "fixtures/test_story.json"

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

var testTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Test Template</title>
</head>
<body>
  <section class="page">
    <h1>{{.Title}}</h1>
    {{range .Paragraphs}}
      <p>{{.}}</p>
    {{end}}
    {{if .Options}}
      <ul>
      {{range .Options}}
        <li><a href="/{{.Chapter}}">{{.Text}}</a></li>
      {{end}}
      </ul>
    {{else}}
      <h3>The End</h3>
    {{end}}
  </section>
</body>
</html>`

var testPathFn = func(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)
	if path == "/story" || path == "/story/" {
		path = "/story/intro"
	}
	return path[len("/story/"):]
}

func TestNewHandler(t *testing.T) {
	story := createTestStory()
	customTmpl := template.Must(template.New("").Parse(testTemplate))
	templateOpt := WithTemplate(customTmpl)
	pathFnOpt := WithPathFunc(testPathFn)

	tests := []struct {
		desc               string
		opts               []HandlerOption
		wantCustomTemplate bool
		wantCustomPathFn   bool
	}{
		{
			"without HandlerOpts",
			nil,
			false,
			false,
		},
		{
			"with template HandlerOpt",
			[]HandlerOption{templateOpt},
			true,
			false,
		},
		{
			"with path func HanlderOpt",
			[]HandlerOption{pathFnOpt},
			false,
			true,
		},
		{
			"with both HandlerOpts",
			[]HandlerOption{templateOpt, pathFnOpt},
			true,
			true,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			// Set defaultPathFn to nil for easy comparison
			oldPathFn := defaultPathFn
			defaultPathFn = nil

			// Configure handler
			h := NewHandler(story, test.opts...)
			concreteH := h.(*handler)

			// Compare received template with default
			gotTmpl := concreteH.t
			hasDefaultTmpl := gotTmpl == defaultTmpl
			if hasDefaultTmpl == test.wantCustomTemplate {
				t.Errorf(
					"NewHandler(story, [customTmpl: %t, customPathFn: %t]): customTemplate == %t, want %t",
					test.wantCustomTemplate, test.wantCustomPathFn, !hasDefaultTmpl, test.wantCustomTemplate,
				)
			}

			// Compare receive pathFn with default
			gotPathFn := concreteH.pathFn
			hasDefaultPathFn := gotPathFn == nil
			if hasDefaultPathFn == test.wantCustomPathFn {
				t.Errorf(
					"NewHandler(story, [customTmpl: %t, customPathFn: %t]): customPathFn == %t, want %t",
					test.wantCustomTemplate, test.wantCustomPathFn, !hasDefaultPathFn, test.wantCustomPathFn,
				)
			}

			// Clean up
			defaultPathFn = oldPathFn
		})
	}
}

func TestHandlerServeHTTP(t *testing.T) {
	story := createTestStory()
	h := NewHandler(story)

	tests := []struct {
		path              string
		wantStatusCode    int
		wantBodyToInclude string
	}{
		{"/C1", 200, "This is Chapter 1."},
		{"/not-found", 404, ""},
	}

	for _, test := range tests {
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, test.path, nil)
		h.ServeHTTP(rec, req)
		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != test.wantStatusCode {
			t.Errorf("%s: got status %d %s, want %d %s", test.path,
				res.StatusCode, http.StatusText(res.StatusCode),
				test.wantStatusCode, http.StatusText(test.wantStatusCode),
			)
		}

		if test.wantBodyToInclude != "" {
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Errorf("%s: %v", test.path, err)
			}
			if !bytes.Contains(body, []byte(test.wantBodyToInclude)) {
				t.Errorf("%s: want body to contain %q; got body %q",
					test.path, test.wantBodyToInclude, body)
			}
		}
	}
}

func createTestStory() Story {
	f, err := os.Open(storyFixture)
	if err != nil {
		panic(err.Error())
	}
	defer f.Close()
	story, err := JSONStory(f)
	if err != nil {
		panic(err.Error())
	}
	return story
}
