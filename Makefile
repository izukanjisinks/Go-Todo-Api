# Makefile for Go project
DB_SERVER=localhost
DB_PORT=1433
DB_NAME=Todos


# Database migrations
migrate-up:
	migrate -path migrations -database "sqlserver://$(DB_SERVER):$(DB_PORT)?database=$(DB_NAME)&trusted_connection=yes" up

migrate-down:
	migrate -path migrations -database "sqlserver://$(DB_SERVER):$(DB_PORT)?database=$(DB_NAME)&trusted_connection=yes" down

migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)

migrate-force:
	migrate -path migrations -database "sqlserver://$(DB_SERVER):$(DB_PORT)?database=$(DB_NAME)&trusted_connection=yes" force $(version)

# The name of the compiled binary (executable file)
BINARY_NAME=todo-api

# -----------------------------
# build: Compiles the Go source code into a binary file
# -----------------------------
build:
	go build -o $(BINARY_NAME)

# -----------------------------
# run: Runs your Go project directly (without building a separate binary)
# -----------------------------
run:
	go run ./cmd/api

# -----------------------------
# install: Installs your Go program into the system's Go bin directory
# -----------------------------
install:
	go install

# -----------------------------
# clean: Cleans up build files and removes the compiled binary
# -----------------------------
clean:
	go clean
	rm -f $(BINARY_NAME)

# -----------------------------
# Docker commands
# -----------------------------

# Build Docker image
docker-build:
	docker build -t todo-api:latest .

# Start services (creates containers if they don't exist)
docker-up:
	docker-compose up

# Start services in background (detached mode)
docker-up-d:
	docker-compose up -d

# Stop services (keeps containers)
docker-stop:
	docker-compose stop

# Stop and remove containers
docker-down:
	docker-compose down

# View logs
docker-logs:
	docker-compose logs -f

# Rebuild and restart services
docker-rebuild:
	docker-compose up --build

# List running services
docker-ps:
	docker-compose ps

# Execute shell command in running container
docker-shell:
	docker-compose exec api sh

# Full restart (rebuild image and restart containers)
docker-restart: docker-build docker-down docker-up-d

# Remove all containers, networks, and volumes
docker-clean:
	docker-compose down -v
	docker rmi todo-api:latest
