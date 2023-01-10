package preview_url

import (
	"strings"
	"testing"
)

func TestGetLinkPreviewItems(t *testing.T) {
	uri := "https://doziesiky.com"
	maxRedirect := 10

	scraper := NewScraper(uri, maxRedirect)

	doc, err := scraper.GetLinkPreviewItems()

	if err != nil {
		t.Errorf("Failed to get link preview items: %s", err)
	}

	if doc == nil {
		t.Error("Failed to get link preview items: result is nil")
	}

	if doc.Preview.Link != uri {
		t.Errorf("Expected link to be %s, got %s", uri, doc.Preview.Link)
	}
}

func TestParseDocument(t *testing.T) {
	htmlContent := []byte(`<html>
        <head>
            <meta name="icon" content="http://www.example.com/icon.png">
            <meta name="name" content="Example">
            <meta name="title" content="Example Title">
            <meta name="description" content="Example Description">
            <link rel="icon" href="http://www.example.com/favicon.ico">
            <img src="http://www.example.com/image1.png">
            <img src="http://www.example.com/image2.png">
        </head>
        <body>
            <h1>Example Page</h1>
        </body>
    </html>`)
	doc := &Document{}

	scraper := NewScraper("", 5)
	err := scraper.ParseDocument(htmlContent, doc)

	if err != nil {
		t.Errorf("Failed to parse document: %s", err)
	}

	if doc.Preview.Icon != "http://www.example.com/favicon.ico" {
		t.Errorf("Expected icon to be %s, got %s", "http://www.example.com/icon.png", doc.Preview.Icon)
	}

	if doc.Preview.Name != "Example" {
		t.Errorf("Expected name to be %s, got %s", "Example", doc.Preview.Name)
	}

	if doc.Preview.Title != "Example Title" {
		t.Errorf("Expected title to be %s, got %s", "Example Title", doc.Preview.Title)
	}

	if doc.Preview.Description != "Example Description" {
		t.Errorf("Expected description to be %s, got %s", "Example Description", doc.Preview.Description)
	}

	if len(doc.Preview.Images) != 2 {
		t.Errorf("Expected 2 images, got %d", len(doc.Preview.Images))
	}

	if doc.Preview.Images[0] != "http://www.example.com/image1.png" {
		t.Errorf("Expected first image to be %s, got %s", "http://www.example.com/image1.png", doc.Preview.Images[0])
	}

	if doc.Preview.Images[1] != "http://www.example.com/image2.png" {
		t.Errorf("Expected second image to be %s, got %s", "http://www.example.com/image2.png", doc.Preview.Images[1])
	}
}

func TestToFragmentUrl(t *testing.T) {
	uri := "http://www.example.com#!param1=value1&param2=value2"
	maxRedirect := 10

	scraper := NewScraper(uri, maxRedirect)

	err := scraper.toFragmentUrl()

	if err != nil {
		t.Errorf("Failed to convert to fragment URL: %s", err)
	}

	if !strings.Contains(scraper.EscapedFragmentUrl.String(), "_escaped_fragment_=param1%3Dvalue1%26param2%3Dvalue2") {
		t.Errorf("Expected URL to contain %s, got %s", "_escaped_fragment_=param1%3Dvalue1%26param2%3Dvalue2", scraper.EscapedFragmentUrl.String())
	}
}

func TestAvoidByte(t *testing.T) {
	if !avoidByte(byte(' ')) {
		t.Error("Expected byte ' ' to be avoided")
	}

	if !avoidByte(byte('\n')) {
		t.Error("Expected byte '\n' to be avoided")
	}

	//if !avoidByte(byte('\t')) {
	//	t.Error("Expected byte '\t' to be avoided")
	//}

	if avoidByte(byte('a')) {
		t.Error("Expected byte 'a' to be not avoided")
	}
}

func TestEscapeByte(t *testing.T) {
	if !escapeByte(byte('&')) {
		t.Error("Expected byte '&' to be escaped")
	}

	if !escapeByte(byte('?')) {
		t.Error("Expected byte '?' to be escaped")
	}

	if !escapeByte(byte('=')) {
		t.Error("Expected byte '=' to be escaped")
	}

	if !escapeByte(byte('#')) {
		t.Error("Expected byte '#' to be escaped")
	}

	if !escapeByte(byte('%')) {
		t.Error("Expected byte '%' to be escaped")
	}

	if escapeByte(byte('a')) {
		t.Error("Expected byte 'a' to be not escaped")
	}
}
