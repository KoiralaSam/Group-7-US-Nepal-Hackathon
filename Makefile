# Database URL for golang-migrate (override per environment).
# Example: postgres://USER:PASSWORD@HOST:PORT/DBNAME?sslmode=disable
DATABASE_URL ?= postgres://postgres:postgres@localhost:5432/g7?sslmode=disable
MIGRATIONS_PATH ?= backend/migrations

.PHONY: migrate-up migrate-down migrate-version

migrate-up:
	migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" down 1

migrate-version:
	migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" version
