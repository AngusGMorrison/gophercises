package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/angusgmorrison/gophercises/link"
)

type loc struct {
	Value string `xml:"loc"`
}

type urlset struct {
	URLs  []loc  `xml:"url"`
	Xmlns string `xml:"xmnls,attr"`
}

const xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"

var hostURL *url.URL // match only URLs with the host scheme and domain

func main() {
	URLFlag := flag.String("url", "https://gophercises.com", "the URL of the site to map")
	maxDepth := flag.Int("depth", 3, "the maximum number of links deep to traverse")
	flag.Parse()

	var err error
	hostURL, err = trimURL(*URLFlag)
	if err != nil {
		exit(err.Error())
	}
	pages := bfs(*URLFlag, *maxDepth)
	toXML := urlset{
		URLs:  make([]loc, len(pages)),
		Xmlns: xmlns,
	}
	for i, page := range pages {
		toXML.URLs[i] = loc{page}
	}
	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("", "  ")
	fmt.Print(xml.Header)
	if err = enc.Encode(toXML); err != nil {
		exit(fmt.Sprintf("encoding XML: %v", err))
	}
	fmt.Println()

}

// trimPath returns a *url.URL consisting of only sceheme and host.
func trimURL(URLStr string) (*url.URL, error) {
	parsedURL, err := url.Parse(URLStr)
	if err != nil {
		return nil, fmt.Errorf("parsing %s: %v", URLStr, err)
	}
	schemeHost := &url.URL{
		Scheme: parsedURL.Scheme,
		Host:   parsedURL.Host,
	}
	return schemeHost, nil
}

// Sending searchResults that record their own depth to worklist over a channel results in ~20%
// fewer memory allocations than the traditional approach of rotating a working slice and a next
// slice.
type searchResult struct {
	depth int
	links []string
}

func bfs(URLStr string, maxDepth int) []string {
	seen := make(map[string]struct{})
	worklist := make(chan searchResult)
	sema := make(chan struct{}, 20) // counting semaphore

	go func() { worklist <- searchResult{links: []string{URLStr}} }()

	for pending := 1; pending > 0; pending-- { // run until all goroutines have returned
		results := <-worklist
		for _, link := range results.links {
			if _, ok := seen[link]; !ok {
				seen[link] = struct{}{}
				if results.depth == maxDepth {
					continue // record links found at maxDepth as seen, but do not crawl
				}
				pending++
				go func(URLStr string) {
					sema <- struct{}{}
					worklist <- searchResult{results.depth + 1, get(URLStr)}
					<-sema
				}(link)
			}
		}
	}

	foundLinks := make([]string, 0, len(seen))
	for link := range seen {
		foundLinks = append(foundLinks, link)
	}
	return foundLinks
}

// get fetches a webpage and returns all links to other pages in the host domain.
func get(URLStr string) []string {
	resp, err := http.Get(URLStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "getting %s: %v\n", URLStr, err)
		return []string{}
	}
	foundLinks := filter(hrefs(resp.Body, hostURL), withPrefix(hostURL.String()))
	resp.Body.Close()
	return foundLinks
}

// filter takes a slice of links and retains those for which keepFn returns true.
func filter(links []string, keepFn func(string) bool) []string {
	var ret []string
	for _, l := range links {
		if keepFn(l) {
			ret = append(ret, l)
		}
	}
	return ret
}

// hrefs extracts returns the links from an HTML page, stripping their query strings and expanding
// relative paths.
func hrefs(r io.Reader, base *url.URL) []string {
	links, _ := link.Parse(r)
	var ret []string
	for _, link := range links {
		parsedURL, err := url.Parse(link.Href)
		if err != nil {
			continue
		}
		parsedURL.RawQuery = "" // strip query string to prevent multiple variations of same URL
		parsedURL.Fragment = ""
		resolvedLink := base.ResolveReference(parsedURL)
		ret = append(ret, resolvedLink.String())
	}
	return ret
}

// withPrefix returns a func which determines whether a given URL string has the specified prefix.
func withPrefix(pfx string) func(string) bool {
	return func(link string) bool {
		return strings.HasPrefix(link, pfx)
	}
}

func exit(msg string) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], msg)
	os.Exit(1)
}
