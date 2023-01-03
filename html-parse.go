package wiktionary

import (
	"log"
	"strings"

	"golang.org/x/net/html"
)

func parseHtmlTable(text string) {
	// parse the HTML into a tree
	doc, err := html.Parse(strings.NewReader(text))
	if err != nil {
		log.Fatal(err)
	}
	var f func(*html.Node, int)
	f = func(n *html.Node, depth int) {
		outputData(n, depth)

		// recursive processing
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c, depth+1)
		}
	}
	f(doc, 0)
}

// testing function - delete eventually
func outputData(n *html.Node, depth int) {
	t := strings.Repeat("| ", depth)
	switch n.Type {
	case html.TextNode:
		t += "Text: "
	case html.DocumentNode:
		t += "Document: "
	case html.ElementNode:
		t += "Element: "
	case html.CommentNode:
		t += "Comment: "
	}
	t += n.Data + " "
	for _, attr := range n.Attr {
		t += "[" + attr.Key + " " + attr.Val + "] "
	}
	println(t)
}
