package main

import (
	"hotdeal-tracker/internal/config"
	"hotdeal-tracker/internal/database"
	"hotdeal-tracker/internal/models"
	"log"
)

func main() {
	cfg, err := config.LoadConfig("./config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.NewDatabase(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Testing database operations...")

	// Test: Count products
	var productCount int64
	db.DB.Model(&models.Product{}).Count(&productCount)
	log.Printf("Total products in database: %d", productCount)

	// Test: Get all products
	products, _, err := db.GetAllProducts(10, 0)
	if err != nil {
		log.Printf("Failed to get products: %v", err)
	} else {
		log.Printf("Fetched %d products", len(products))
		for i, p := range products[:min(3, len(products))] {
			log.Printf("%d. Title: %s, Price: %.2f, Platform: %s", i+1, truncate(p.Title, 30), p.Price, p.Platform)
		}
	}

	// Test: Get platforms
	platforms, err := db.GetPlatforms()
	if err != nil {
		log.Printf("Failed to get platforms: %v", err)
	} else {
		log.Printf("Available platforms: %d", len(platforms))
		for _, p := range platforms {
			log.Printf("  - %s (%s): %s", p.Name, p.Code, p.HotURL)
		}
	}

	// Test: Insert test product
	testProduct := models.Product{
		Title:        "Test Product",
		ProductURL:   "https://example.com/test",
		Platform:     "test",
		PlatformID:   "test-123",
		Price:        99.99,
		OriginalPrice: 129.99,
		Rating:       4.5,
		ReviewCount:  100,
		SalesCount:   1000,
		IsHot:        true,
		TrendingScore: 95.0,
	}

	err = db.CreateProduct(&testProduct)
	if err != nil {
		log.Printf("Failed to create test product: %v", err)
	} else {
		log.Printf("Created test product with ID: %d", testProduct.ID)
	}

	// Verify insertion
	db.DB.Model(&models.Product{}).Count(&productCount)
	log.Printf("Total products after insert: %d", productCount)

	log.Println("Database test completed!")
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