This code is a Go implementation of a link preview scraper. Given a URL, it can fetch the contents of the web page at that URL and extract some metadata (icon, name, title, description, images, and link) to create a DocumentPreview struct. This struct can be embedded in a Document struct which also contains the HTML contents of the page in a bytes.Buffer.

The process of creating the preview is initiated by calling the GetLinkPreviewItems(uri string, maxRedirect int) function, which takes in a string containing a URL and an integer specifying the maximum number of redirects to follow. This function then creates a new Scraper struct, sets the Url and MaxRedirect fields, and then calls the GetLinkPreviewItems() method on the struct which performs the actual scraping.

The Scraper struct has a Url field, which is a pointer to a url.URL struct representing the URL to be scraped, and a MaxRedirect field, which stores the maximum number of redirects to follow.

The Scraper struct also has a EscapedFragmentUrl field which is used to store a modified URL that includes a "escaped_fragment" parameter, which some websites use to provide alternative content for crawlers, if the original URL contains a "#!" fragment. This field is used in the getUrl() method which returns either the EscapedFragmentUrl or Url field depending on if a EscapedFragmentUrl is present.

The Scraper struct has a toFragmentUrl() method that modifies the URL by replacing the "#!" fragment with a "escaped_fragment" parameter, when the URL contains a "#!" fragment, in order to provide alternative content for the scraper.

The Scraper struct has a getDocument() method that is used to fetch the HTML contents of the page at the specified URL, either the Url or EscapedFragmentUrl fields.

Overall, this code defines several structs, type and a few methods which when combined allows to fetch the webpage, HTML content, meta data of webpage and parse the webpage, extract metadata and create DocumentPreview struct.

type Scraper struct:

Url *url.URL: The URL of the web page to scrape.
MaxRedirect int: The maximum number of redirects allowed.
client *http.Client : The http.Client struct to be used when sending the request to get the link preview.
NewScraper(uri string, maxRedirect int) *Scraper : This function creates a new instance of the Scraper struct. It takes the uri as the parameter and parse it to *url.URL struct. It also creates a new client struct with CheckRedirect option set and sets the client in the Scraper struct.

func (scraper *Scraper) GetLinkPreviewItems() (*Document, error): This function retrieves the link preview items for the specified URL. If the URL contains a "#!", it converts it to escaped fragment URL, It creates a new GET request to the URL and sends it using the http client struct. After it receives the response, it creates a new buffer, copies the response body to it and passes it with the Document struct to parseDocument().

func (scraper *Scraper) toFragmentUrl() error : This function converts the "#!" in the URL to the escaped fragment URL format.
func (scraper *Scraper) parseDocument(htmlContent []byte, doc *Document) error: This function parses the HTML content passed as a byte slice and extracts the required metadata from it like icon, name, title, description, images, and link. It uses the recursive function to traverse the nodes in the HTML tree, check for the required elements and extract the metadata. It sets the link field of the DocumentPreview struct to the URL of the page, so that the final output will have the complete link of the page.

func avoidByte(b byte) bool: This function takes a byte as input and checks if the given byte should be avoided, the provided example uses bytes for spaces, newlines, etc.

func escapeByte(b byte) bool: This function takes a byte as input and checks if the given byte should be escaped, the provided example uses bytes '&', '?', '=', '#', '%' as bytes to escape, but you can adjust it according to your requirements.

type Document struct: This struct represents the final output of the scrape, it has two fields:

Body bytes.Buffer : The buffer which holds the HTML content.
Preview DocumentPreview : A struct that holds the preview data of the page like icon, name, title, description, images, and link
type DocumentPreview struct: This struct holds the preview data of the page like icon, name, title, description, images, and link.

Overall this implementation uses the Go's builtin http package and url package to send a GET request to the provided URL and parse the response to extract metadata. The package also uses the Go's built-in html package to parse the HTML content to extract the required metadata.
