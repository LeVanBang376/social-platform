APP_NAME=supermarket_server_app

MIGRATIONS_PATH=./migrations

# Load .env file if it exists
ifneq (,$(wildcard .env))
	include .env
	export
endif

.PHONY: run build test tidy migrate-up migrate-down migrate-create migrate-version migrate-force docker-up docker-down seed swag

## ===== GO APP =====

run:
	go run cmd/api/main.go

build:
	go build -o $(APP_NAME) .

tidy:
	go mod tidy

test:
	go test ./...

## ===== MIGRATION =====

migrate-create:
	migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq $(name)

migrate-up:
	migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" down 1

migrate-version:
	migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" version

migrate-force:
	migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" force $(version)

## ==== SEED =====
seed:
	go run ./cmd/seed

## ===== DOCKER =====

docker-up:
	docker compose up -d

docker-down:
	docker compose down

## ===== SWAGGER =====
swag:
	swag init -g cmd/api/main.go --parseDependency --parseInternal

## ===== RESET =====
# reset:
# 	docker compose -f docker-compose.database.yml down -v
# 	docker compose -f docker-compose.database.yml up -d --build
# 	@timeout /t 3 /nobreak > nul
# 	$(MAKE) migrate-up
# 	$(MAKE) seed
reset:
	dropdb -U postgres --if-exists smallsupermarket_database
	createdb -U postgres smallsupermarket_database
	$(MAKE) migrate-up
	$(MAKE) seed