package main

import (
	"hotdeal-tracker/pkg/parser"
	"io/ioutil"
	"log"
)

func main() {
	data, err := ioutil.ReadFile("ebay.html")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Testing parser with saved HTML (length: %d bytes)", len(data))

	p := parser.GetParser("ebay")
	products, err := p.ParseProducts(data, "ebay")
	if err != nil {
		log.Fatalf("Failed to parse products: %v", err)
	}

	log.Printf("Found %d products total", len(products))

	// Show products with prices > 0
	log.Println("\n--- Products with prices ---")
	count := 0
	for i, product := range products {
		if product.Price > 0 {
			log.Printf("\nProduct %d:", i+1)
			log.Printf("  Title: %s", truncate(product.Title, 50))
			log.Printf("  Price: %.2f (%s)", product.Price, product.Currency)
			log.Printf("  URL: %s", truncate(product.ProductURL, 80))
			count++
			if count >= 5 {
				break
			}
		}
	}

	log.Printf("\nTotal products with price > 0: %d", count)
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}