# Simple Makefile for a Go project

# Build the application
all: build

build:
	@echo "Building..."
	@templ generate
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
	@if command -v air > /dev/null; then \
	    air; \
	    echo "Watching...";\
	else \
	    read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
	    if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
	        go install github.com/air-verse/air@latest; \
	        air; \
	        echo "Watching...";\
	    else \
	        echo "You chose not to install air. Exiting..."; \
	        exit 1; \
	    fi; \
	fi

.PHONY: all build run test clean tailwind sqlc
