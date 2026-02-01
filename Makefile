# Simple Makefile for a Go project

# Build the application
all: build

build:
	@echo "Building..."
	@go run github.com/a-h/templ/cmd/templ@latest generate
	@npx tailwindcss -i ./static/assets/css/input.css -o ./static/assets/css/output.css
	@go build -o main cmd/web/main.go

# Run the application
run:
	@go run cmd/web/main.go



# Test the application
test:
	@echo "Testing..."
	@go test ./... -v

sqlc:
	sqlc generate

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

tailwind:
	@npx tailwindcss -i static/assets/css/input.css -o static/assets/css/output.css
	
# Live Reload
watch:
	@go run github.com/air-verse/air@latest

.PHONY: all build run test clean tailwind sqlc
