package preview_url

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewScraper(t *testing.T) {
	scraper, err := NewScraper("http://example.com", 10)
	if err != nil {
		t.Fatalf("Failed to create scraper: %s", err)
	}
	if scraper.BaseURL.String() != "http://example.com" {
		t.Fatalf("Unexpected BaseURL, got: %s, want: http://example.com", scraper.BaseURL.String())
	}
	if scraper.MaxRedirects != 10 {
		t.Fatalf("Unexpected MaxRedirects, got: %d, want: 10", scraper.MaxRedirects)
	}
}

func TestGetPreviewMetadata(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<html><head><title>Test Page</title></head><body><img src='test.png'></body></html>"))
	}))
	defer ts.Close()

	scraper, _ := NewScraper(ts.URL, 10)
	doc, err := scraper.GetPreviewMetadata()
	if err != nil {
		t.Fatalf("Failed to get preview metadata: %s", err)
	}
	if doc.Preview.Title != "Test Page" {
		t.Fatalf("Unexpected title, got: %s, want: Test Page", doc.Preview.Title)
	}
	if len(doc.Preview.Images) != 1 || doc.Preview.Images[0] != "test.png" {
		t.Fatalf("Unexpected images, got: %v, want: [test.png]", doc.Preview.Images)
	}
}

func TestMaxRedirects(t *testing.T) {
	redirections := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if redirections < 5 {
			http.Redirect(w, r, "/redirect", http.StatusFound)
			redirections++
			return
		}
		w.Write([]byte("<html><head><title>Test Page</title></head><body><img src='test.png'></body></html>"))
	}))
	defer ts.Close()

	scraper, _ := NewScraper(ts.URL, 3)
	_, err := scraper.GetPreviewMetadata()
	if err == nil || err.Error() != "exceeded max redirects: 4" {
		t.Fatalf("Expected to exceed max redirects, got: %v", err)
	}
}

func TestEscapedFragmentURL(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("_escaped_fragment_") != "fragment" {
			t.Fatalf("Expected to receive escaped fragment, got: %s", r.URL.Query().Get("_escaped_fragment_"))
		}
		w.Write([]byte("<html><head><title>Test Page</title></head><body><img src='test.png'></body></html>"))
	}))
	defer ts.Close()

	scraper, _ := NewScraper(ts.URL+"#!fragment", 10)
	doc, err := scraper.GetPreviewMetadata()
	if err != nil {
		t.Fatalf("Failed to get preview metadata: %s", err)
	}
	if doc.Preview.Title != "Test Page" {
		t.Fatalf("Unexpected title, got: %s, want: Test Page", doc.Preview.Title)
	}
	if len(doc.Preview.Images) != 1 || doc.Preview.Images[0] != "test.png" {
		t.Fatalf("Unexpected images, got: %v, want: [test.png]", doc.Preview.Images)
	}
}
