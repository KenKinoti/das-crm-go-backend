.PHONY: build run dev test clean install-dev

# Build the application
build:
	go build -o bin/server cmd/server/main.go

# Run the application
run: build
	./bin/server

# Run in development mode with hot reload (requires air)
dev:
	air

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/
	rm -rf tmp/

# Install development dependencies
install-dev:
	go install github.com/air-verse/air@latest

# Generate Swagger documentation
swagger:
	swag init -g cmd/server/main.go

# Install all Go dependencies
deps:
	go mod tidy
	go mod download

# Run without hot reload
run-simple:
	go run cmd/server/main.go
