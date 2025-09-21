.PHONY: build test clean run-server run-cli install-frontend build-frontend docker-build docker-run

# Go binary name
BINARY_NAME=csv2json

# Build the Go binary
build:
	go build -o $(BINARY_NAME) .

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	go clean
	rm -f $(BINARY_NAME)
	rm -rf frontend/build
	rm -rf frontend/node_modules

# Run as server
run-server: build
	./$(BINARY_NAME) -server

# Run CLI with sample data
run-cli: build
	echo "name,age,city\nJohn,30,NYC\nJane,25,LA" > sample.csv
	./$(BINARY_NAME) -i sample.csv
	rm sample.csv

# Install frontend dependencies
install-frontend:
	cd frontend && npm install

# Build frontend
build-frontend: install-frontend
	cd frontend && npm run build

# Build everything
build-all: build-frontend build

# Docker build
docker-build:
	docker build -t csv2json .

# Docker run
docker-run: docker-build
	docker run -p 8080:8080 csv2json

# Development setup
dev-setup: install-frontend
	go mod tidy

# Format code
fmt:
	go fmt ./...
	cd frontend && npm run format 2>/dev/null || true

# Lint code
lint:
	golangci-lint run
	cd frontend && npm run lint 2>/dev/null || true
