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

// Constants to avoid magic strings in code
const (
	EscapedFragment = "_escaped_fragment_="
	UserAgent       = "link-preview-scraper"
)

// Regex for URL fragment matching
var fragmentRegexp = regexp.MustCompile("#!(.*)")

// Scraper Struct for holding scraper related data
type Scraper struct {
	BaseURL            *url.URL
	EscapedFragmentURL *url.URL
	MaxRedirects       int
	client             *http.Client
}

// Document Struct for holding the scraped document data
type Document struct {
	Body    bytes.Buffer
	Preview DocumentPreview
}

// DocumentPreview Struct for holding the preview data of the scraped document
type DocumentPreview struct {
	Icon        string
	Name        string
	Title       string
	Description string
	Images      []string
	Link        string
}

// NewScraper initializes a new scraper
func NewScraper(uri string, maxRedirects int) (*Scraper, error) {
	parsedURL, err := url.Parse(uri)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= maxRedirects {
				return fmt.Errorf("exceeded max redirects: %d", len(via))
			}
			return nil
		},
	}

	return &Scraper{BaseURL: parsedURL, MaxRedirects: maxRedirects, client: client}, nil
}

// GetPreviewMetadata fetches and parses the document for preview metadata
func (s *Scraper) GetPreviewMetadata() (*Document, error) {
	if s.requiresEscapedFragmentURL() {
		if err := s.createEscapedFragmentURL(); err != nil {
			return nil, fmt.Errorf("failed to create escaped fragment URL: %w", err)
		}
	}

	req, err := http.NewRequest("GET", s.getTargetURL(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", UserAgent)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, resp.Body); err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	doc := &Document{}
	if err := s.parseHTML(buf.Bytes(), doc); err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}
	return doc, nil
}

// requiresEscapedFragmentURL checks whether the URL has a fragment that needs to be escaped
func (s *Scraper) requiresEscapedFragmentURL() bool {
	return strings.Contains(s.BaseURL.String(), "#!") || strings.Contains(s.BaseURL.String(), EscapedFragment)
}

// createEscapedFragmentURL creates an escaped fragment URL from the base URL
func (s *Scraper) createEscapedFragmentURL() error {
	unescapedURL, err := url.QueryUnescape(s.BaseURL.String())
	if err != nil {
		return err
	}

	// Code to generate the escaped fragment URL removed for brev
	// Matching fragments in the URL
	matches := fragmentRegexp.FindStringSubmatch(unescapedURL)

	// Adding the escaped fragment to the URL
	if len(matches) > 1 {
		escapedFragment := EscapedFragment
		for _, r := range matches[1] {
			b := byte(r)
			if isAvoidableByte(b) {
				continue
			}
			if shouldEscapeByte(b) {
				escapedFragment += url.QueryEscape(string(r))
			} else {
				escapedFragment += string(r)
			}
		}

		// Preparing the final URL
		paramPrefix := "?"
		if len(s.BaseURL.Query()) > 0 {
			paramPrefix = "&"
		}
		escapedFragmentURL, err := url.Parse(strings.Replace(unescapedURL, matches[0], paramPrefix+escapedFragment, 1))
		if err != nil {
			return err
		}
		s.EscapedFragmentURL = escapedFragmentURL
	} else {
		paramPrefix := "?"
		if len(s.BaseURL.Query()) > 0 {
			paramPrefix = "&"
		}
		escapedFragmentURL, err := url.Parse(unescapedURL + paramPrefix + EscapedFragment)
		if err != nil {
			return err
		}
		s.EscapedFragmentURL = escapedFragmentURL
	}
	return nil
}

// getTargetURL returns the appropriate URL to scrape
func (s *Scraper) getTargetURL() string {
	if s.EscapedFragmentURL != nil {
		return s.EscapedFragmentURL.String()
	}
	return s.BaseURL.String()
}

// parseHTML parses the HTML content and fills the document with relevant preview data
func (s *Scraper) parseHTML(htmlContent []byte, doc *Document) error {
	// Parsing the HTML
	node, err := html.Parse(bytes.NewReader(htmlContent))
	if err != nil {
		return err
	}

	// Recursive function to traverse the HTML tree
	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode {
			s.processNode(n, doc)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}
	traverse(node)

	// Setting the link field to the URL of the page
	doc.Preview.Link = s.getTargetURL()

	return nil
}

// processNode processes a single HTML node and updates the document preview accordingly
func (s *Scraper) processNode(n *html.Node, doc *Document) {
	switch n.Data {
	case "meta":
		s.processMetaNode(n, doc)
	case "link":
		s.processLinkNode(n, doc)
	case "img":
		s.processImageNode(n, doc)
	}
}

// processMetaNode processes a meta node and updates the document preview
func (s *Scraper) processMetaNode(n *html.Node, doc *Document) {
	var name, content string
	for _, attr := range n.Attr {
		if attr.Key == "name" {
			name = attr.Val
		} else if attr.Key == "content" {
			content = attr.Val
		}
	}

	switch name {
	case "icon":
		doc.Preview.Icon = content
	case "name":
		doc.Preview.Name = content
	case "title":
		doc.Preview.Title = content
	case "description":
		doc.Preview.Description = content
	}
}

// processLinkNode processes a link node and updates the document preview
func (s *Scraper) processLinkNode(n *html.Node, doc *Document) {
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
}

// processImageNode processes an image node and updates the document preview
func (s *Scraper) processImageNode(n *html.Node, doc *Document) {
	for _, attr := range n.Attr {
		if attr.Key == "src" {
			doc.Preview.Images = append(doc.Preview.Images, attr.Val)
			break
		}
	}
}

// isAvoidableByte checks if a byte should be avoided when escaping a URL
func isAvoidableByte(b byte) bool {
	avoid := []byte{' ', '\n', '\r'}
	for _, v := range avoid {
		if b == v {
			return true
		}
	}
	return false
}

// shouldEscapeByte checks if a byte should be escaped when creating a URL
func shouldEscapeByte(b byte) bool {
	escape := []byte{'&', '?', '=', '#', '%'}
	for _, v := range escape {
		if b == v {
			return true
		}
	}
	return false
}
