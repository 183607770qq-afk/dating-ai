package main

import (
	"flag"
	"hotdeal-tracker/internal/config"
	"hotdeal-tracker/pkg/httpclient"
	"hotdeal-tracker/pkg/parser"
	"log"
	"os"
)

func main() {
	url := flag.String("url", "https://www.amazon.com/gp/bestsellers", "URL to test")
	platform := flag.String("platform", "amazon", "Platform code")
	output := flag.String("output", "", "Output file for HTML")
	flag.Parse()

	log.Printf("Testing URL: %s", *url)
	log.Printf("Platform: %s", *platform)

	cfg := &config.CrawlerConfig{
		UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		Timeout:   30,
	}

	client := httpclient.NewHTTPClient(cfg)

	log.Println("Fetching HTML...")
	html, err := client.Get(*url, nil)
	if err != nil {
		log.Fatalf("Failed to fetch URL: %v", err)
	}

	log.Printf("HTML length: %d bytes", len(html))

	if *output != "" {
		os.WriteFile(*output, html, 0644)
		log.Printf("Saved HTML to %s", *output)
	}

	log.Println("Parsing products...")
	p := parser.GetParser(*platform)
	products, err := p.ParseProducts(html, *platform)
	if err != nil {
		log.Fatalf("Failed to parse products: %v", err)
	}

	log.Printf("Found %d products", len(products))

	if len(products) > 0 {
		log.Println("--- First 3 products ---")
		for i, product := range products[:min(3, len(products))] {
			log.Printf("%d. Title: %s", i+1, truncate(product.Title, 50))
			log.Printf("   Price: $%.2f", product.Price)
			log.Printf("   URL: %s", truncate(product.ProductURL, 80))
			log.Println()
		}
	} else {
		log.Println("No products found. Trying GenericParser...")
		gp := &parser.GenericParser{}
		products, err = gp.ParseProducts(html, *platform)
		if err != nil {
			log.Printf("GenericParser error: %v", err)
		}
		log.Printf("GenericParser found %d products", len(products))
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}