package utils

import (
	"io"
	"net/http"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func ParseTitleByURL(url string) (string, bool, error) {
	resp, err := http.DefaultClient.Get(url)
	if err != nil {
		return "", false, err
	}
	return parseTitle(resp.Body)
}

func parseTitle(reader io.Reader) (string, bool, error) {
	node, err := html.Parse(reader)
	if err != nil {
		return "", false, err
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if c.DataAtom != atom.Html {
			continue
		}
		for c2 := c.FirstChild; c2 != nil; c2 = c2.NextSibling {
			if c2.DataAtom != atom.Head {
				continue
			}
			for c3 := c2.FirstChild; c3 != nil; c3 = c3.NextSibling {
				if c3.DataAtom != atom.Title {
					continue
				}
				for c4 := c3.FirstChild; c4 != nil; c4 = c4.NextSibling {
					if c4.Type != html.TextNode {
						continue
					}
					return c4.Data, true, nil
				}
			}
		}
	}
	return "", false, nil
}
