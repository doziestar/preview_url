# Preview_url

This package provides functionality to scrape and parse metadata from a given URL. The metadata can be used to generate a link preview, similar to what you might see in a messaging app when a link is shared. The package provides an easy-to-use interface to initialize a new scraper, fetch, and parse the document for preview metadata.

## Installation
```
go get github.com/doziestar/preview_url
```

## Usage

To use this package, import it into your Go code, create a new scraper with the URL and maximum number of redirects, then call GetPreviewMetadata() to retrieve the document's preview data.
    
    ```go
    import (
        "fmt"
        "github.com/doziestar/preview_url"
    )

    scraper, err := preview_url.NewScraper(""http://example.com", 10)
    if err != nil {
    log.Fatal(err)
    }
    
    doc, err := scraper.GetPreviewMetadata()
    if err != nil {
    log.Fatal(err)
    }
    
    fmt.Println(doc.Preview.Title)
    fmt.Println(doc.Preview.Description)
    fmt.Println(doc.Preview.Images)
    ```

### Dependencies

This package relies on the `net/http` and `net/url` packages from the standard library, as well as `golang.org/x/net/html` for parsing HTML.

### Testing
Unit tests for this package can be written using Go's built-in testing framework. These tests can check that the scraper is correctly initializing, that it is handling redirects properly, and that it is correctly parsing the HTML and extracting the preview metadata.

### Structs 

The scraper uses the following structs to store the scraper results and the preview metadata.

```go
type Scraper struct {
	BaseURL            *url.URL
	EscapedFragmentURL *url.URL
	MaxRedirects       int
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
```

### Functions and Methods
### `NewScraper`
This function initializes a new Scraper.
    
    ```go
    func NewScraper(uri string, maxRedirects int) (*Scraper, error)
    ```

- uri - The URL to be scraped.
- maxRedirects - The maximum number of redirects allowed.

GetPreviewMetadata
This method fetches and parses the document for preview metadata.
    
    ```go
    func (s *Scraper) GetPreviewMetadata() (*Document, error)
    ```

- Returns a pointer to a Document struct and an error.
- The Document struct contains the scraped document's body and preview metadata.
- The error is nil if the document was successfully scraped and parsed.

### `getEscapedFragmentURL`
This method returns the URL with the escaped fragment.
    
    ```go
    func (s *Scraper) getEscapedFragmentURL() (*url.URL, error)
    ```

- Returns a pointer to a url.URL struct and an error.
- The url.URL struct contains the URL with the escaped fragment.
- The error is nil if the URL was successfully escaped.

### `requiresEscapedFragmentURL`
This method returns true if the URL requires an escaped fragment.
    
    ```go
    func (s *Scraper) requiresEscapedFragmentURL() bool
    ```

- Returns a boolean value.
- The boolean value is true if the URL requires an escaped fragment.

### `createEscapedFragmentURL`
This method returns the URL with the escaped fragment.
    
    ```go
    func (s *Scraper) createEscapedFragmentURL() error
    ```

- Returns an error.
- The error is nil if the URL was successfully escaped.

### `getTargetURL`
This method returns the target URL.
    
    ```go
    func (s *Scraper) getTargetURL() string
    ```

- Returns a string.
- The string is the target URL.

### `parseHTML`

This method parses the HTML document for preview metadata.
    
    ```go
    func (s *Scraper) parseHTML(htmlContent []byte, doc *Document) error

    ```

- Returns an error.
- The error is nil if the HTML was successfully parsed.

