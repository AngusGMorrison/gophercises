package link

import (
	"bytes"
	"io"
	"strings"

	"golang.org/x/net/html"
)

// Link represents a link in an HTML document (<a href="...">)
type Link struct {
	Href string
	Text string
}

// Parse will take in an HTML document and will return a slice of links parsed from it.
func Parse(r io.Reader) ([]Link, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	nodes := linkNodes(doc)
	links := make([]Link, 0, len(nodes))
	for _, n := range nodes {
		links = append(links, buildLink(n))
	}
	return links, nil
}

func linkNodes(n *html.Node) []*html.Node {
	var ret []*html.Node
	var dfs func(n *html.Node)
	dfs = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			ret = append(ret, n)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			dfs(c)
		}
	}
	dfs(n)
	return ret
}

func buildLink(n *html.Node) Link {
	var ret Link
	for _, attr := range n.Attr {
		if attr.Key == "href" {
			ret.Href = attr.Val
			break
		}
	}
	ret.Text = text(n)
	return ret
}

func text(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	if n.Type != html.ElementNode {
		return ""
	}

	buf := bytes.NewBuffer([]byte{})
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		buf.WriteString(text(c))
	}
	str := buf.String()
	// Clean up irregular whitespace
	return strings.Join(strings.Fields(str), " ")
}
