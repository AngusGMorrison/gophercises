package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/angusgmorrison/gophercises/cyoa"
)

func main() {
	filename := flag.String("file", "gopher.json", "the JSON file with the CYOA story")
	port := flag.Int("port", 3000, "the port to start the CYOA web application on")
  flag.Parse()
  fmt.Printf("Using the story in %s.\n", *filename)

  // Open the JSON file and parse the story in it.
	f, err := os.Open(*filename)
	if err != nil {
		exit(err.Error())
	}
	story, err := cyoa.JSONStory(f)
	if err != nil {
		exit(fmt.Sprintf("parsing %s: %v", *filename, err))
	}

  // Create our customer CYOA story handler
	tmpl := template.Must(template.New("").Parse(storyTmpl))
	handler := cyoa.NewHandler(
		story,
		cyoa.WithTemplate(tmpl),
		cyoa.WithPathFunc(pathFn),
	)

  // Create a ServeMux to route our requests
	mux := http.NewServeMux()
	mux.Handle("/story/", handler)
	log.Printf("Starting the server at: %d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), mux))
}

func pathFn(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)
	if path == "/story" || path == "/story/" {
		path = "/story/intro"
	}
	return path[len("/story/"):]
}

var storyTmpl = `
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
    <ul>
    {{range .Options}}
      <li><a href="/story/{{.Chapter}}">{{.Text}}</a></li>
    {{end}}
	</ul>
  </section>
</body>
</html>`

func exit(msg string) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], msg)
	os.Exit(1)
}
