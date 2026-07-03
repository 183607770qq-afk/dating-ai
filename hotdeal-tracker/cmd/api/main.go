package main

import (
	"flag"
	"fmt"
	"hotdeal-tracker/internal/api"
	"hotdeal-tracker/internal/config"
	"hotdeal-tracker/internal/database"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	configPath := flag.String("config", "./config.yaml", "Path to configuration file")
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

	handler := api.NewHandler(db)

	router := api.SetupRouter(handler, cfg.Server.Mode)

	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Starting API server on %s", addr)

	go func() {
		if err := router.Run(addr); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
}
