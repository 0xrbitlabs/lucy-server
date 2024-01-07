.PHONY: setup-db stop-db delete-db reboot-db

build:
	@go build -o server ./main.go

run-dev:
	@go run main.go

install-deps:
	@go mod download

setup-db:
	@docker run --name db -e POSTGRES_PASSWORD=pwd -p 5432:5432 -itd postgres:latest

stop-db:
	@docker stop db

delete-db:
	@docker rm db

reboot-db: stop-db delete-db setup-db

migrate:
	@docker cp ./schema.sql db:/tmp/schema.sql
	@docker exec -it db psql -U postgres -d bizapp -f /tmp/schema.sql
	@echo "Migration successful"

build-image:
	@docker build -t server:latest .

push-image:
	@docker tag server:latest thewisepigeon/visio:latest
	@docker push thewisepigeon/server:latest

test:
	@go test -v ./...
