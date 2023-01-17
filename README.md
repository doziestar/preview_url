This code is a Go implementation of a link preview scraper. Given a URL, it can fetch the contents of the web page at that URL and extract some metadata (icon, name, title, description, images, and link) to create a DocumentPreview struct. This struct can be embedded in a Document struct which also contains the HTML contents of the page in a bytes.Buffer.

The process of creating the preview is initiated by calling the GetLinkPreviewItems(uri string, maxRedirect int) function, which takes in a string containing a URL and an integer specifying the maximum number of redirects to follow. This function then creates a new Scraper struct, sets the Url and MaxRedirect fields, and then calls the GetLinkPreviewItems() method on the struct which performs the actual scraping.

The Scraper struct has a Url field, which is a pointer to a url.URL struct representing the URL to be scraped, and a MaxRedirect field, which stores the maximum number of redirects to follow.

The Scraper struct also has a EscapedFragmentUrl field which is used to store a modified URL that includes a "escaped_fragment" parameter, which some websites use to provide alternative content for crawlers, if the original URL contains a "#!" fragment. This field is used in the getUrl() method which returns either the EscapedFragmentUrl or Url field depending on if a EscapedFragmentUrl is present.

The Scraper struct has a toFragmentUrl() method that modifies the URL by replacing the "#!" fragment with a "escaped_fragment" parameter, when the URL contains a "#!" fragment, in order to provide alternative content for the scraper.

The Scraper struct has a getDocument() method that is used to fetch the HTML contents of the page at the specified URL, either the Url or EscapedFragmentUrl fields.
