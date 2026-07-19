# Create the initial postgres container
# NOTE: Removed sudo so Docker Desktop can see the container
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