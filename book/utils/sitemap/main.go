package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// XML structures for sitemap
type URLSet struct {
	XMLName xml.Name `xml:"urlset"`
	Xmlns   string   `xml:"xmlns,attr"`
	URLs    []URL    `xml:"url"`
}

type URL struct {
	Loc        string `xml:"loc"`
	LastMod    string `xml:"lastmod,omitempty"`
	ChangeFreq string `xml:"changefreq,omitempty"`
	Priority   string `xml:"priority,omitempty"`
}

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <book-dir> <base-url>\n", os.Args[0])
		os.Exit(1)
	}

	bookDir := os.Args[1]
	baseURL := strings.TrimSuffix(os.Args[2], "/")

	urlset := URLSet{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  []URL{},
	}

	// Walk through the book directory to find HTML files
	err := filepath.Walk(bookDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip non-HTML files
		if !strings.HasSuffix(path, ".html") {
			return nil
		}

		// Skip print.html if it exists
		if strings.HasSuffix(path, "print.html") {
			return nil
		}

		// Get relative path from book directory
		relPath, err := filepath.Rel(bookDir, path)
		if err != nil {
			return err
		}

		// Convert Windows paths to URL paths
		urlPath := filepath.ToSlash(relPath)

		// Create full URL
		fullURL := baseURL + "/" + urlPath

		// Get file modification time
		lastMod := info.ModTime().Format("2006-01-02T15:04:05-07:00")

		// Determine priority and change frequency
		priority := "0.5"
		changeFreq := "weekly"

		// Higher priority for important pages
		if urlPath == "index.html" {
			priority = "1.0"
			changeFreq = "daily"
		} else if strings.Contains(urlPath, "quick-start") ||
			strings.Contains(urlPath, "getting-started") ||
			strings.Contains(urlPath, "introduction") {
			priority = "0.8"
			changeFreq = "weekly"
		} else if strings.Contains(urlPath, "tutorial") {
			priority = "0.7"
		}

		// Add URL to sitemap
		urlset.URLs = append(urlset.URLs, URL{
			Loc:        fullURL,
			LastMod:    lastMod,
			ChangeFreq: changeFreq,
			Priority:   priority,
		})

		return nil
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error walking directory: %v\n", err)
		os.Exit(1)
	}

	// Generate XML
	output, err := xml.MarshalIndent(urlset, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating XML: %v\n", err)
		os.Exit(1)
	}

	// Write XML header and content
	fmt.Printf("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n%s\n", output)
}
