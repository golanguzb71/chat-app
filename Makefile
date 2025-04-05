APP_NAME = chat-app
BINARY_NAME = chat-app

# Default
.PHONY: all
all: run

# Run the application
.PHONY: run
run:
	go run cmd/main.go

# Build the binary
.PHONY: build
build:
	go build -o bin/$(BINARY_NAME) cmd/main.go

# Tidy dependencies
.PHONY: tidy
tidy:
	go mod tidy

# Run Docker container
.PHONY: docker-up
docker-up:
	docker-compose up --build

# Stop Docker container
.PHONY: docker-down
docker-down:
	docker-compose down

.PHONY: create-migration
create-migration:
ifndef NAME
	$(error NAME is not set. Usage: make create-migration NAME=create_users_table)
endif
	migrate create -ext sql -dir migrations -seq $(NAME)