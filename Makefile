include .envrc

# ========================================= #
# Helpers
# ========================================== #

## help: prints this help message
.PHONY: help
help:
	@echo "Usage: "
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ":" | sed -e "s/^/ /"

.PHONY: confirm
confirm:
	@echo -n "Are you sure? [y/N] " && read ans && [ $${ans:-N} = y ]

# ========================================= #
# Development
# ========================================== #

## run/api: run the cmd/api application
.PHONY: run/api
run/api:
	go build -o bin/api ./cmd/api/
	exec bin/api -db-dsn=${PFAPI-DB-DSN}

## db/migrations/new name=$1: create a new database migration
.PHONY: db/migrations/new
db/migrations/new:
	@echo "Creating migration files for ${name}..."
	migrate create -seq -ext=.sql -dir=./migrations ${name}
	
## db/migrations/up: apply all database migrations
.PHONY: db/migrations/up
db/migrations/up: confirm
	@echo "Running migrations..."
	migrate -path ./migrations -database=${PFAPI-DB-DSN} up

# ========================================= #
# Quality Control
# ========================================== #

## audit: tidy dependencies and format, vet and test all code
.PHONY: audit
audit:
	@echo "Tidying and veryfing module dependencies..."
	go mod tidy
	go mod verify
	@echo "Formatting code..."
	go fmt ./...
	@echo "Vetting code..."
	go vet ./...
	staticcheck ./...

# ========================================= #
# Build
# ========================================== #

## build/api: build the cmd/api application
.PHONY: build/api
build/api:
	@echo "Building cmd/api..."
	go build -ldflags="-s" -o=./bin/api ./cmd/api
