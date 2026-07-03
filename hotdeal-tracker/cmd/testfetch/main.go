package main

import (
	"hotdeal-tracker/internal/config"
	"hotdeal-tracker/pkg/httpclient"
	"log"
	"os"
)

func main() {
	url := "https://www.ebay.com/b/Best-Selling-Best-Sellers-on-eBay/6000"

	cfg := &config.CrawlerConfig{
		UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		Timeout:   30,
	}

	client := httpclient.NewHTTPClient(cfg)

	log.Println("Fetching HTML...")
	html, err := client.Get(url, nil)
	if err != nil {
		log.Fatalf("Failed to fetch URL: %v", err)
	}

	log.Printf("HTML length: %d bytes", len(html))

	err = os.WriteFile("ebay_new.html", html, 0644)
	if err != nil {
		log.Fatalf("Failed to write file: %v", err)
	}

	log.Println("Saved to ebay_new.html")
	log.Println("First 500 bytes:")
	log.Println(string(html[:min(500, len(html))]))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}