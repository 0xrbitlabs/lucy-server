include .env

build:
	@go build -o ./tmp/main .

install-deps:
	@go mod download

dev:
	@go build -o ./tmp/main cmd/main.go
	@./tmp/main

setup-db:
	@docker exec -it $(DB_NAME) psql -U postgres -c 'CREATE DATABASE lucy;'

migrate:
	@docker cp ./schema.sql $(DB_NAME):/tmp/schema.sql
	@docker exec -it $(DB_NAME) psql -U postgres -d lucy -f /tmp/schema.sql
	@echo "Migration successful"

test:
	@go test -v ./...
