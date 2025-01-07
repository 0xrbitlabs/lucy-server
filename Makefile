include .env

DB_PORT ?= 5432
BINARY=lucy.out

build:
	go build -o $(BINARY) .

up-db: down-db
	@docker run --name=$(DB_CONTAINER_NAME) -e POSTGRES_PASSWORD=$(DB_PASSWORD) -e POSTGRES_DB=$(DB_NAME) -itd -p 5432:$(DB_PORT) postgres:latest

down-db:
	@docker stop $(DB_CONTAINER_NAME) && docker rm $(DB_CONTAINER_NAME) || true

reset-db:
	@docker exec -it $(DB_CONTAINER_NAME) psql -U postgres -d $(DB_NAME) -c "DROP SCHEMA public CASCADE;"
	@docker exec -it $(DB_CONTAINER_NAME) psql -U postgres -d $(DB_NAME) -c "CREATE SCHEMA public;"
	$(MAKE) migrate

into-db:
	@docker exec -it $(DB_CONTAINER_NAME) bash

migrate:
	@docker cp ./schema.sql $(DB_CONTAINER_NAME):/tmp/schema.sql
	@docker exec -it $(DB_CONTAINER_NAME) psql -U postgres -d $(DB_NAME) -f /tmp/schema.sql

seed:
	@docker cp ./seed.sql $(DB_CONTAINER_NAME):/tmp/seed.sql
	@docker exec -it $(DB_CONTAINER_NAME) psql -U postgres -d $(DB_NAME) -f /tmp/seed.sql

ngrok:
	@ngrok http --url=trusty-serval-master.ngrok-free.app 8080
