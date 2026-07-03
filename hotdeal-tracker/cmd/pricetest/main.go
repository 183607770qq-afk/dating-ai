package main

import (
	"hotdeal-tracker/internal/config"
	"hotdeal-tracker/pkg/httpclient"
	"hotdeal-tracker/pkg/parser"
	"log"
)

func main() {
	url := "https://www.ebay.com/b/Best-Selling-Best-Sellers-on-eBay/6000"
	platform := "ebay"

	log.Printf("Testing price extraction from %s", url)

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

	log.Println("Parsing products...")
	p := parser.GetParser(platform)
	products, err := p.ParseProducts(html, platform)
	if err != nil {
		log.Fatalf("Failed to parse products: %v", err)
	}

	log.Printf("Found %d products total", len(products))

	// Show first 3 CSS products (likely categories)
	log.Println("\n--- First 3 CSS products (categories) ---")
	for i, product := range products[:min(3, len(products))] {
		log.Printf("\nProduct %d:", i+1)
		log.Printf("  Title: %s", product.Title)
		log.Printf("  Price: %.2f (%s)", product.Price, product.Currency)
	}

	// Show last 3 products (should be from JSON)
	log.Println("\n--- Last 3 products (from JSON) ---")
	startIdx := max(0, len(products)-3)
	for i := startIdx; i < len(products); i++ {
		product := products[i]
		log.Printf("\nProduct %d:", i+1)
		log.Printf("  Title: %s", product.Title)
		log.Printf("  Price: %.2f (%s)", product.Price, product.Currency)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}