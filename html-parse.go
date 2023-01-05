package wiktionary

import (
	"log"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

func parseInflectionTable(part *PartOfSpeech, text string) {
	// parse the HTML into a tree
	doc, err := html.Parse(strings.NewReader(text))
	if err != nil {
		log.Fatal(err)
	}
	var f func(*html.Node, int)
	f = func(n *html.Node, depth int) {
		// if this is an element node which has the attribute key "***form-of"
		// then capture this specific form as an extended part
		// NB will be similar to 1&#124;s&#124;pres&#124;act&#124;ind-form-of
		// or 1|s|pres|act|ind-form-of once converted
		if n.Type == html.ElementNode {
			for _, attr := range n.Attr {
				if strings.Contains(attr.Key, "form-of") && len(attr.Key) > 8 {
					// convert the &#124; to |
					form := strings.ReplaceAll(attr.Key, "&#124;", "|")
					// strip off the "-form-of" part, it's superfluous
					form = form[:len(form)-8]
					// get the matching text
					word := getElementText(n)
					// add it to the parts list
					// first check to see if that part already exists
					if _, ok := part.Attributes[form]; ok {
						v := 2
						for {
							// if an alternative version with this version number exists
							if _, ok := part.Attributes[form+"|alt"+strconv.Itoa(v)]; ok {
								// increment the version number and try again
								v += 1
							} else {
								// otherwise add the part with this version number
								part.Attributes[form+"|alt"+strconv.Itoa(v)] = word
								break
							}
						}

					} else {
						// if it doesn't exist then just add it
						part.Attributes[form] = word
					}
				}
			}
		}

		// recursive processing
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c, depth+1)
		}
	}
	f(doc, 0)
}

func getElementText(n *html.Node) string {
	// get the enclosed text for an HTML element
	text := ""
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		switch c.Type {
		case html.TextNode:
			text += c.Data
		case html.ElementNode:
			text += getElementText(c)
		}
	}
	return text
}
