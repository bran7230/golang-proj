# Create the initial postgres container
# NOTE: Removed sudo so Docker Desktop can see the container
-include .env
export
# Define paths
BINARY_NAME := golang-proj
CMD_DIR := ./cmd/golang-proj
BIN_DIR := bin
TESTING_DIR := ./internal/tests

# Build target
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BIN_DIR)
# Build the package located in CMD_DIR and output to BIN_DIR/BINARY_NAME
	go vet ./internal/*
	go vet ./cmd/*
	go build -o $(BIN_DIR)/$(BINARY_NAME) $(CMD_DIR)

# Run target
run: build
	./$(BIN_DIR)/$(BINARY_NAME)

# Clean target
clean:
	@if [ -d "$(BIN_DIR)" ]; then \
		echo "Cleaning $(BIN_DIR)..."; \
		rm -rf $(BIN_DIR)/*; \
		echo "Bin cleaned!"; \
	else \
		echo "Directory '$(BIN_DIR)' does not exist, nothing to clean."; \
	fi

test:
	go vet ./internal/*
	go vet ./cmd/*
	go test $(TESTING_DIR) -v
create-docker-db:
	docker run --name $(DOCKER_DB_NAME) \
	   -e POSTGRES_PASSWORD=$(DOCKER_PASS) \
	   -e POSTGRES_DB=$(POSTGRESS_DB_NAME)\
	   -p $(DOCKER_PORT):5432 \
	   -d postgres:latest

# Remove the docker container
clean-docker-db:
	docker rm -f $(DOCKER_DB_NAME)

# Check container status
check-docker-health:
	docker ps

# Connect to postgres interactively
# Updated to use the environment variable instead of a hardcoded name
connect-to-db:
	docker exec -it $(DOCKER_DB_NAME) psql -U postgres -d $(POSTGRESS_DB_NAME)

# Stop the database container
stop-docker-db:
	docker stop $(DOCKER_DB_NAME)

# Start the database container
start-docker-db:
	docker start $(DOCKER_DB_NAME)