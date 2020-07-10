package cyoa

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
)

func init() {
	defaultTmpl = template.Must(template.New("").Parse(defaultTmplStr))
}

var defaultTmpl *template.Template

var defaultTmplStr = `
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Choose Your Own Adventure</title>
  <style>
      body {
        font-family: helvetica, arial;
      }
      h1 {
        text-align:center;
        position:relative;
      }
      .page {
        width: 80%;
        max-width: 500px;
        margin: auto;
        margin-top: 40px;
        margin-bottom: 40px;
        padding: 80px;
        background: #FFFCF6;
        border: 1px solid #eee;
        box-shadow: 0 10px 6px -6px #777;
      }
      ul {
        border-top: 1px dotted #ccc;
        padding: 10px 0 0 0;
        -webkit-padding-start: 0;
      }
      li {
        padding-top: 10px;
      }
      a,
      a:visited {
        text-decoration: none;
        color: #6295b5;
      }
      a:active,
      a:hover {
        color: #7792a2;
      }
      p {
        text-indent: 1em;
      }
    </style>
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

// A Story is the CYOA superstructure, mapping chapter title strings to the chapter contents.
type Story map[string]Chapter

// A Chapter holds the details of a single step within the CYOA journey.
type Chapter struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []Option `json:"options"`
}

// An Option contains a reference to another chapter to be displayed if the user selects the
// corresponding text.
type Option struct {
	Text    string `json:"text"`
	Chapter string `json:"chapter"`
}

// JSONStory decodes a story from input JSON.
func JSONStory(r io.Reader) (Story, error) {
	d := json.NewDecoder(r)
	var story Story
	if err := d.Decode(&story); err != nil {
		return nil, err
	}
	return story, nil
}

type handler struct {
	s      Story
	t      *template.Template
	pathFn func(r *http.Request) string
}

var defaultPathFn = func(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)
	if path == "" || path == "/" {
		path = "/intro"
	}
	return path[1:]
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := h.pathFn(r)
	if chapter, ok := h.s[path]; ok {
		if err := h.t.Execute(w, chapter); err != nil {
			log.Printf("%v", err)
			http.Error(w, "Something went wrong...", http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, "Chapter not found.", http.StatusNotFound)
}

// HandlerOption returns a function that configures a handler when called.
type HandlerOption func(h *handler)

// WithTemplate allows users to specify a custom template for the handler to use in responses.
func WithTemplate(t *template.Template) HandlerOption {
	return func(h *handler) {
		h.t = t
	}
}

// WithPathFunc allows users to specify an alternative URL structure to the handler. By default,
// the hanlder will match all paths.
func WithPathFunc(fn func(r *http.Request) string) HandlerOption {
	return func(h *handler) {
		h.pathFn = fn
	}
}

// NewHandler configures a new handler according to the story and HanlderOptions supplied by the
// user, and returns a pointer that satisfies the http.Handler interface.
func NewHandler(s Story, opts ...HandlerOption) http.Handler {
	h := &handler{s, defaultTmpl, defaultPathFn}
	for _, opt := range opts {
		opt(h)
	}
	return h
}
