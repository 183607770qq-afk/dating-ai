.PHONY: all build test clean run-api run-scraper migrate

all: build

build:
	@echo "Building API server..."
	@go build -o bin/api cmd/api/main.go
	@echo "Building scraper..."
	@go build -o bin/scraper cmd/scraper/main.go
	@echo "Build complete!"

test:
	go test ./...

clean:
	rm -rf bin/

run-api:
	go run cmd/api/main.go -config ./config.yaml

run-scraper:
	go run cmd/scraper/main.go -config ./config.yaml

run-scraper-once:
	go run cmd/scraper/main.go -config ./config.yaml -once

crawl-keyword:
	go run cmd/scraper/main.go -config ./config.yaml -keyword "$(KEYWORD)" -platform "$(PLATFORM)"

deps:
	go mod download
	go mod tidy

lint:
	golangci-lint run

migrate:
	psql -U postgres -d hotdeal_tracker -f migrations/001_init.sql

docker-build:
	docker build -t hotdeal-tracker .

docker-run:
	docker run -p 8080:8080 hotdeal-tracker

help:
	@echo "Available commands:"
	@echo "  make build           - Build all binaries"
	@echo "  make test            - Run tests"
	@echo "  make clean           - Clean build artifacts"
	@echo "  make run-api         - Run API server"
	@echo "  make run-scraper     - Run scraper service"
	@echo "  make run-scraper-once - Run single crawl cycle"
	@echo "  make crawl-keyword   - Crawl specific keyword (KEYWORD=xxx PLATFORM=xxx)"
	@echo "  make deps            - Download dependencies"
	@echo "  make lint            - Run linter"
	@echo "  make migrate         - Run database migrations"
