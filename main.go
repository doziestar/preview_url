package preview_url

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

/**
This code is a Go implementation of a link preview scraper. Given a URL, it can fetch the contents of the web page at that URL and extract some metadata (icon, name, title, description, images, and link) to create a DocumentPreview struct. This struct can be embedded in a Document struct which also contains the HTML contents of the page in a bytes.Buffer.

The process of creating the preview is initiated by calling the GetLinkPreviewItems(uri string, maxRedirect int) function, which takes in a string containing a URL and an integer specifying the maximum number of redirects to follow. This function then creates a new Scraper struct, sets the Url and MaxRedirect fields, and then calls the GetLinkPreviewItems() method on the struct which performs the actual scraping.

The Scraper struct has a Url field, which is a pointer to a url.URL struct representing the URL to be scraped, and a MaxRedirect field, which stores the maximum number of redirects to follow.

The Scraper struct also has a EscapedFragmentUrl field which is used to store a modified URL that includes a "escaped_fragment" parameter, which some websites use to provide alternative content for crawlers, if the original URL contains a "#!" fragment. This field is used in the getUrl() method which returns either the EscapedFragmentUrl or Url field depending on if a EscapedFragmentUrl is present.

The Scraper struct has a toFragmentUrl() method that modifies the URL by replacing the "#!" fragment with a "escaped_fragment" parameter, when the URL contains a "#!" fragment, in order to provide alternative content for the scraper.

The Scraper struct has a getDocument() method that is used to fetch the HTML contents of the page at the specified URL, either the Url or EscapedFragmentUrl fields.

Overall, this code defines several structs, type and a few methods which when combined allows to fetch the webpage, HTML content, meta data of webpage and parse the webpage, extract metadata and create DocumentPreview struct.
*/

var (
	EscapedFragment string = "_escaped_fragment_="
	fragmentRegexp         = regexp.MustCompile("#!(.*)")
)

type Scraper struct {
	Url                *url.URL
	EscapedFragmentUrl *url.URL
	MaxRedirect        int
	client             *http.Client
}

type Document struct {
	Body    bytes.Buffer
	Preview DocumentPreview
}

type DocumentPreview struct {
	Icon        string
	Name        string
	Title       string
	Description string
	Images      []string
	Link        string
}

func NewScraper(uri string, maxRedirect int) *Scraper {
	u, err := url.Parse(uri)
	if err != nil {
		return nil
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= maxRedirect {
				return fmt.Errorf("too many redirects")
			}
			return nil
		},
	}

	return &Scraper{Url: u, MaxRedirect: maxRedirect, client: client}
}

func (scraper *Scraper) GetLinkPreviewItems() (*Document, error) {
	if strings.Contains(scraper.Url.String(), "#!") {
		if err := scraper.toFragmentUrl(); err != nil {
			return nil, err
		}
	}
	if strings.Contains(scraper.Url.String(), EscapedFragment) {
		scraper.EscapedFragmentUrl = scraper.Url
	}

	req, err := http.NewRequest("GET", scraper.getUrl(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "link-preview-scraper")

	resp, err := scraper.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, resp.Body)
	if err != nil {
		return nil, err
	}

	doc := &Document{}
	err = scraper.parseDocument(buf.Bytes(), doc)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func (scraper *Scraper) toFragmentUrl() error {
	unescapedurl, err := url.QueryUnescape(scraper.Url.String())
	if err != nil {
		return err
	}
	matches := fragmentRegexp.FindStringSubmatch(unescapedurl)
	if len(matches) > 1 {
		escapedFragment := EscapedFragment
		for _, r := range matches[1] {
			b := byte(r)
			if avoidByte(b) {
				continue
			}
			if escapeByte(b) {
				escapedFragment += url.QueryEscape(string(r))
			} else {
				escapedFragment += string(r)
			}
		}

		p := "?"
		if len(scraper.Url.Query()) > 0 {
			p = "&"
		}
		fragmentUrl, err := url.Parse(strings.Replace(unescapedurl, matches[0], p+escapedFragment, 1))
		if err != nil {
			return err
		}
		scraper.EscapedFragmentUrl = fragmentUrl
	} else {
		p := "?"
		if len(scraper.Url.Query()) > 0 {
			p = "&"
		}
		fragmentUrl, err := url.Parse(unescapedurl + p + EscapedFragment)
		if err != nil {
			return err
		}
		scraper.EscapedFragmentUrl = fragmentUrl
	}
	return nil
}

func (scraper *Scraper) getUrl() string {
	if scraper.EscapedFragmentUrl != nil {
		return scraper.EscapedFragmentUrl.String()
	}
	return scraper.Url.String()
}

func (scraper *Scraper) parseDocument(htmlContent []byte, doc *Document) error {
	// Parsing the HTML
	node, err := html.Parse(bytes.NewReader(htmlContent))
	if err != nil {
		return err
	}

	// Use a recursive function to traverse the nodes in the HTML tree
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "meta" {
			var name, content string
			for _, attr := range n.Attr {
				if attr.Key == "name" {
					name = attr.Val
				} else if attr.Key == "content" {
					content = attr.Val
				}
			}

			if name == "icon" {
				doc.Preview.Icon = content
			} else if name == "name" {
				doc.Preview.Name = content
			} else if name == "title" {
				doc.Preview.Title = content
			} else if name == "description" {
				doc.Preview.Description = content
			}
		} else if n.Type == html.ElementNode && n.Data == "link" {
			var rel, href string
			for _, attr := range n.Attr {
				if attr.Key == "rel" {
					rel = attr.Val
				} else if attr.Key == "href" {
					href = attr.Val
				}
			}
			if rel == "icon" {
				doc.Preview.Icon = href
			}
		} else if n.Type == html.ElementNode && n.Data == "img" {
			for _, attr := range n.Attr {
				if attr.Key == "src" {
					doc.Preview.Images = append(doc.Preview.Images, attr.Val)
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(node)

	// Setting the link field to the URL of the page
	doc.Preview.Link = scraper.getUrl()

	return nil
}

func avoidByte(b byte) bool {
	// List of bytes to avoid, e.g. spaces, newlines, etc.
	avoid := []byte{' ', '\n', '\r'}
	for _, v := range avoid {
		if b == v {
			return true
		}
	}
	return false
}

func escapeByte(b byte) bool {
	// List of bytes to escape
	escape := []byte{'&', '?', '=', '#', '%'}
	for _, v := range escape {
		if b == v {
			return true
		}
	}
	return false
}
