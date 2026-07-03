package main

import (
	"hotdeal-tracker/internal/config"
	"hotdeal-tracker/pkg/httpclient"
	"log"
	"os"
	"strings"
)

func main() {
	url := "https://www.ebay.com/b/Best-Selling-Best-Sellers-on-eBay/6000"

	cfg := &config.CrawlerConfig{
		UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		Timeout:   30,
	}

	client := httpclient.NewHTTPClient(cfg)

	log.Println("Fetching HTML...")
	htmlBytes, err := client.Get(url, nil)
	if err != nil {
		log.Fatalf("Failed to fetch URL: %v", err)
	}

	html := string(htmlBytes)
	log.Printf("HTML length: %d bytes", len(html))

	// Save to file for analysis
	err = os.WriteFile("ebay.html", htmlBytes, 0644)
	if err != nil {
		log.Fatalf("Failed to write file: %v", err)
	}

	log.Println("HTML saved to ebay.html")
	log.Println("Searching for price patterns...")

	// Look for price patterns in the HTML
	if strings.Contains(html, "$") {
		log.Println("Found $ symbol")
	}
	if strings.Contains(strings.ToLower(html), "price") {
		log.Println("Found 'price' text")
	}
	if strings.Contains(html, "USD") {
		log.Println("Found 'USD' text")
	}
}