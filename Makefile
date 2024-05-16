BINARY_NAME=main
PORT=8080

build:
	@echo "Building $(BINARY_NAME)..."
	@go build  -o ./bin/$(BINARY_NAME) ./app
	@chmod +x ./bin/$(BINARY_NAME)

run: build
	@echo "Running $(BINARY_NAME) on port $(PORT)..."
	@./bin/$(BINARY_NAME)

clean:
	@echo "Cleaning up..."
	@go clean
	@rm -f ./bin/$(BINARY_NAME)

.PHONY: build run clean


test:
	@echo "Running tests..."
	@go test ./... -v