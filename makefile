# Define paths
BINARY_NAME := golang-proj
CMD_DIR := ./cmd/golang-proj
BIN_DIR := bin

# Build target
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BIN_DIR)
# Build the package located in CMD_DIR and output to BIN_DIR/BINARY_NAME
	go build -o $(BIN_DIR)/$(BINARY_NAME) $(CMD_DIR)

# Run target
run: build
	./$(BIN_DIR)/$(BINARY_NAME)

# Clean target
clean:
	rm -rf $(BIN_DIR)/*
	rm -rf docs

.PHONY: build run clean