package main

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	
	"github.com/gocolly/colly/v2"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run script.go <website_url>")
		os.Exit(1)
	}

	// The website URL to scrape is taken from the first command-line argument.
	websiteURL := os.Args[1]

	// Parse the URL and use the Host name as the file base.
	parsedURL, err := url.Parse(websiteURL)
	if err != nil {
		fmt.Printf("Could not parse web URL: %v\n", err)
		os.Exit(1)
	}
	host := strings.Replace(parsedURL.Hostname(), "www.", "", -1) // Remove 'www.' if it exists.
	outputDir := "output"
	filename := fmt.Sprintf("%s.html", host)

	// Ensure the output directory exists.
	err = os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		fmt.Printf("Failed to create output directory: %v\n", err)
		os.Exit(1)
	}

	// Create the output file within the "output" directory.
	filepath := filepath.Join(outputDir, filename)
	file, err := os.Create(filepath)
	if err != nil {
		fmt.Printf("Failed to create file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// Create a new collector.
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"),
	)

	// On every HTML element which has a body tag, perform the following function.
	c.OnHTML("body", func(e *colly.HTMLElement) {
		htmlContent, err := e.DOM.Html()
		if err != nil {
			fmt.Printf("Failed to get HTML content: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("HTML body found, saving to %s.\n", filepath)
		// Save to file.
		_, err = file.WriteString(htmlContent)
		if err != nil {
			fmt.Printf("Failed to write to file: %v\n", err)
			os.Exit(1)
		}
	})

	// Handle any errors.
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Error:", err.Error())
	})

	// Start scraping the website.
	err = c.Visit(websiteURL)
	if err != nil {
		fmt.Printf("Failed to scrape %s: %v\n", websiteURL, err)
		os.Exit(1)
	}
}