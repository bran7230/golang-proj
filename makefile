-include .env
export

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

# Create the inital container
create-docker-db:
	sudo docker run --name $(DOCKER_DB_NAME) -e MYSQL_ROOT_PASSWORD=$(DOCKER_PASS) -p $(DOCKER_PORT):3306 -d mysql:latest

clean-docker-db:
	sudo docker rm -f $(DOCKER_DB_NAME)

check-docker-health:
	sudo docker ps

connect-to-db:
	sudo docker exec -it $(DOCKER_DB_NAME) mysql -u root -p

stop-docker-db:
	sudo docker stop $(DOCKER_DB_NAME)

start-docker-db:
	sudo docker start $(DOCKER_DB_NAME)

.PHONY: build run clean launch-docker-db clean-docker-db check-docker-health connect-to-db