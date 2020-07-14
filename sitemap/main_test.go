package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
)

const testHTML = `
<body>
	<a href="https://test.com/page1">Page 1</a>
	<a href="https://test.com">Page 2</a>
	<a href="/test.com/page3?query=hasone">Page 3</a>
</body>
`

func TestTrimURL(t *testing.T) {
	const wantURL = "https://www.test.com"
	tests := []string{
		"https://www.test.com",
		"https://www.test.com/",
		"https://www.test.com/path",
		"https://www.test.com/path/longer",
		"https://www.test.com/path#fragment",
		"https://www.test.com/path?query=test",
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("trimURL(%q)", test), func(t *testing.T) {
			trimmed, err := trimURL(test)
			if err != nil {
				t.Fatal(err.Error())
			}
			got := trimmed.String()
			if got != wantURL {
				t.Fatalf("got %s, want %s", got, wantURL)
			}
		})
	}
}

const (
	siteFixture          = "fixtures/test_site/"
	defaultServerAddress = "127.0.0.1:49855"
	homePage             = "http://" + defaultServerAddress
	contactPage          = homePage + "/contact"
	aboutPage            = homePage + "/about"
	termsPage            = aboutPage + "/terms"
	cookiesPage          = termsPage + "/cookies"
)

func TestBfs(t *testing.T) {
	server := testServer(defaultServerAddress)
	defer server.Close()
	hostURL, _ = url.Parse(homePage)

	tests := []struct {
		depth     int
		wantLinks []string
	}{
		{1, []string{homePage, aboutPage, contactPage}},
		{2, []string{homePage, aboutPage, contactPage, termsPage}},
		{3, []string{homePage, aboutPage, contactPage, termsPage, cookiesPage}},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("bfs(%q, %d)", homePage, test.depth), func(t *testing.T) {
			results := make(map[string]int)
			gotLinks := bfs(homePage, test.depth)
			for _, link := range gotLinks {
				results[link]++
			}
			for _, link := range test.wantLinks {
				results[link]--
			}
			for _, v := range results {
				if v != 0 {
					t.Fatalf("got %v, want %v", gotLinks, test.wantLinks)
				}
			}
		})
	}
}

func testServer(address string) *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlerFunc("index.html"))
	mux.HandleFunc("/about", handlerFunc("about.html"))
	mux.HandleFunc("/contact", handlerFunc("contact.html"))
	mux.HandleFunc("/about/terms", handlerFunc("terms.html"))
	server := httptest.NewUnstartedServer(mux)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}
	server.Listener = listener
	server.Start()
	return server
}

func handlerFunc(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		renderFile(path, w)
	}
}

func renderFile(path string, w http.ResponseWriter) {
	f, err := os.Open(siteFixture + path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	w.Write(bytes)
}
