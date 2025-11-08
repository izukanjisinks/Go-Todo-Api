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
