package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"net/url"
	"golang.org/x/net/html"
	"hash/fnv" // Import the fnv package for hash generation
)

func main() {
	url := "https://cs272-0304-f23.github.io/tests/top10/"
	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error fetching the URL: %v\n", err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		fmt.Printf("HTTP request failed with status code: %d\n", response.StatusCode)
		return
	}

	doc, err := html.Parse(response.Body)
	if err != nil {
		fmt.Printf("Error parsing HTML: %v\n", err)
		return
	}

	outputFolder := "output_links"
	err = os.MkdirAll(outputFolder, os.ModePerm)
	if err != nil {
		fmt.Printf("Error creating output folder: %v\n", err)
		return
	}
	extractAndProcessLinks(doc, outputFolder)
}

func extractAndProcessLinks(n *html.Node, outputFolder string) {
	if n == nil {
		return
	}

	if n.Type == html.ElementNode && n.Data == "a" {
		link := extractLink(n)
		processLink(link, outputFolder)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		extractAndProcessLinks(c, outputFolder)
	}
}

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

func processLink(link, outputFolder string) {
	if link == "" {
		return
	}

	filename := filepath.Join(outputFolder, url.QueryEscape(link))
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating file %s: %v\n", filename, err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(link)
	if err != nil {
		fmt.Printf("Error writing to file %s: %v\n", filename, err)
		return
	}

	fmt.Printf("Link %s saved to %s\n", link, filename)
}

func split(bookText string, hashFilePath string) (chapterMap map[uint32]string, err error) {
	chapterDelimiter := "CHAPTER"
	chapters := strings.Split(bookText, chapterDelimiter)

	chapterMap = make(map[uint32]string)

	// Load the hash data from the provided file
	hashData, err := loadHashFile(hashFilePath)
	if err != nil {
		return nil, err
	}

	for i := 1; i < len(chapters); i++ {
		title := strings.TrimSpace(chapters[i-1])
		content := strings.TrimSpace(chapters[i])

		// Generate a hash for the chapter title
		hash := hashString(title)

		// Use the hash as the key in the map
		chapterMap[hash] = content
	}

	return chapterMap, nil
}

func loadHashFile(hashFilePath string) (map[string]uint32, error) {
	// Implement your logic to load the hash file here
	// This function should return a map of chapter titles to their corresponding hash values
	// You can use any file reading method to load the hash data from the file.
	// For simplicity, I'm assuming a map of chapter titles to uint32 hash values.

	// Initialize an empty hash data map
	hashData := make(map[string]uint32)

	// Read the hash file and populate the hashData map
	// Implement your file reading logic here

	return hashData, nil
}

func hashString(s string) uint32 {
	// Create a new hash object
	h := fnv.New32a()
	// Write the string to the hash object
	h.Write([]byte(s))
	// Return the hash value as a uint32
	return h.Sum32()
}

