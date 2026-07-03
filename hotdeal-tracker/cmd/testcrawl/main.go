package main

import (
	"hotdeal-tracker/internal/config"
	"hotdeal-tracker/internal/crawler"
	"hotdeal-tracker/internal/database"
	"hotdeal-tracker/internal/models"
	"io/ioutil"
	"log"
)

func main() {
	log.Println("Testing crawler with saved eBay HTML...")

	dbConfig := &config.DatabaseConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "root",
		Password: "root",
		Name:     "hotdeal_tracker",
		SSLMode:  "disable",
	}

	db, err := database.NewDatabase(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	crawlerConfig := &config.CrawlerConfig{
		UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		Timeout:   30,
	}

	c := crawler.NewCrawler(db, crawlerConfig)

	data, err := ioutil.ReadFile("ebay.html")
	if err != nil {
		log.Fatalf("Failed to read ebay.html: %v", err)
	}

	log.Printf("Loaded HTML: %d bytes", len(data))

	platform := models.Platform{
		ID:      2,
		Name:    "eBay",
		Code:    "ebay",
		HotURL:  "https://www.ebay.com/b/Best-Selling-Best-Sellers-on-eBay/6000",
		BaseURL: "https://www.ebay.com",
	}

	products, err := c.CrawlHotProductsWithHTML(data, platform)
	if err != nil {
		log.Fatalf("Failed to crawl products: %v", err)
	}

	log.Printf("Successfully crawled %d products!", len(products))

	// Show some products
	log.Println("\n--- Products with prices ---")
	count := 0
	for _, p := range products {
		if p.Price > 0 {
			log.Printf("\nTitle: %s", truncate(p.Title, 60))
			log.Printf("Price: %.2f %s", p.Price, p.Currency)
			log.Printf("URL: %s", p.ProductURL)
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