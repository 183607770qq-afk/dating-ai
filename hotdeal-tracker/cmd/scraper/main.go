package main

import (
	"flag"
	"hotdeal-tracker/internal/config"
	"hotdeal-tracker/internal/crawler"
	"hotdeal-tracker/internal/database"
	"hotdeal-tracker/internal/models"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	configPath := flag.String("config", "./config.yaml", "Path to configuration file")
	once := flag.Bool("once", false, "Run crawler once and exit")
	keyword := flag.String("keyword", "", "Search keyword")
	platform := flag.String("platform", "", "Specific platform to crawl")
	flag.Parse()

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.NewDatabase(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Connected to database successfully")

	c := crawler.NewCrawler(db, &cfg.Crawler)

	if *keyword != "" {
		platformCode := *platform
		if platformCode == "" {
			platformCode = "amazon"
		}

		log.Printf("Crawling keyword '%s' from %s", *keyword, platformCode)
		products, err := c.CrawlKeyword(*keyword, platformCode)
		if err != nil {
			log.Fatalf("Failed to crawl keyword: %v", err)
		}

		log.Printf("Crawled %d products", len(products))
		return
	}

	if *once {
		log.Println("Running single crawl cycle...")
		platforms, err := db.GetPlatforms()
		if err != nil {
			log.Fatalf("Failed to get platforms: %v", err)
		}

		for _, platform := range platforms {
			log.Printf("Crawling platform: %s", platform.Name)
			
			task := models.CrawlTask{
				PlatformID: platform.ID,
				URL:        platform.HotURL,
				Type:       "hot",
				Status:     "pending",
				Priority:   platform.Priority,
			}

			if err := db.CreateCrawlTask(&task); err != nil {
				log.Printf("Failed to create task for %s: %v", platform.Name, err)
				continue
			}
			
			log.Printf("Created task for %s: %s", platform.Name, task.URL)
			
			products, err := c.CrawlHotProducts(platform)
			if err != nil {
				log.Printf("Failed to crawl %s: %v", platform.Name, err)
				task.Status = "failed"
				task.ErrorMsg = err.Error()
			} else {
				log.Printf("Crawled %d products from %s", len(products), platform.Name)
				task.Status = "completed"
				task.ProductCount = len(products)
				completedTime := time.Now()
				task.CompletedAt = &completedTime
			}
			
			if err := db.UpdateCrawlTask(&task); err != nil {
				log.Printf("Failed to update task for %s: %v", platform.Name, err)
			}
		}

		log.Println("Single crawl cycle completed")
		return
	}

	if err := c.Start(); err != nil {
		log.Fatalf("Failed to start crawler: %v", err)
	}

	log.Println("Crawler started successfully")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	ticker := time.NewTicker(6 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-quit:
			log.Println("Shutting down crawler...")
			c.Stop()
			return
		case <-ticker.C:
			log.Println("Running scheduled crawl cycle...")
			if err := c.Start(); err != nil {
				log.Printf("Error in scheduled crawl: %v", err)
			}
		}
	}
}
