# DATABASE_URL from .env when present; override with: make migrate-up DATABASE_URL='...' or export DATABASE_URL.
ifneq (,$(wildcard .env))
DATABASE_URL ?= $(shell sed -n 's/^DATABASE_URL=//p' .env | head -1)
endif
MIGRATIONS_PATH ?= backend/migrations

.PHONY: migrate-up migrate-down migrate-version migrate-create

migrate-up:
	migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" down 1

migrate-version:
	migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" version

# Next sequential pair:  make migrate-create NAME=add_oauth_sessions
migrate-create:
	@if [ -z "$(NAME)" ]; then \
		echo 'Usage: make migrate-create NAME=<snake_case_description>' >&2; \
		exit 1; \
	fi
	migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq $(NAME)
