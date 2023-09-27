package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"net/url"
	//"io"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	// URL of the web page to download
	url := "https://cs272-0304-f23.github.io/tests/top10/"

	// Send an HTTP GET request to the URL
	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error fetching the URL: %v\n", err)
		return
	}
	defer response.Body.Close()

	// Check if the response status code indicates success
	if response.StatusCode != http.StatusOK {
		fmt.Printf("HTTP request failed with status code: %d\n", response.StatusCode)
		return
	}

	// Parse the HTML content
	doc, err := html.Parse(response.Body)
	if err != nil {
		fmt.Printf("Error parsing HTML: %v\n", err)
		return
	}

	// Extract and process links from the HTML
	outputFolder := "output_links"
	err = os.MkdirAll(outputFolder, os.ModePerm)
	if err != nil {
		fmt.Printf("Error creating output folder: %v\n", err)
		return
	}
	extractAndProcessLinks(doc, outputFolder)
}


// extractAndProcessLinks extracts links from the HTML document and processes them
func extractAndProcessLinks(n *html.Node, outputFolder string) {
	if n == nil {
		return
	}

	if n.Type == html.ElementNode && n.Data == "a" {
		// Extract the href attribute (link) from the <a> element
		link := extractLink(n)

		// Process the link (you can replace this with your own logic)
		processLink(link, outputFolder)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		extractAndProcessLinks(c, outputFolder)
	}
}

// extractLink extracts the href attribute (link) from an HTML <a> element
func extractLink(n *html.Node) string {
	var link string
	if n == nil {
		return link
	}

	if n.Type == html.ElementNode && n.Data == "a" {
		for _, attr := range n.Attr {
			if attr.Key == "href" {
				link = attr.Val
				break
			}
		}
	}

	return link
}

// processLink processes a link and stores it to a file
func processLink(link, outputFolder string) {
	if link == "" {
		return
	}

	// Create a file path based on the link and output folder
	filename := filepath.Join(outputFolder, url.QueryEscape(link))

	// Create a file with the link as its content
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating file %s: %v\n", filename, err)
		return
	}
	defer file.Close()

	// Write the link to the file
	_, err = file.WriteString(link)
	if err != nil {
		fmt.Printf("Error writing to file %s: %v\n", filename, err)
		return
	}

	fmt.Printf("Link %s saved to %s\n", link, filename)
}

func split(bookText string) (chapterMap map[string]string, err error) {
	chapterDelimiter := "CHAPTER"
	chapters := strings.Split(bookText, chapterDelimiter)

	chapterMap = make(map[string]string)

	for i := 1; i < len(chapters); i++ {
		title := strings.TrimSpace(chapters[i-1 ])
		content := strings.TrimSpace(chapters[i])
		chapterMap[title] = content
	}

	return chapterMap, nil
}