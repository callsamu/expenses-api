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
	exec bin/api \
		-db-dsn=${EXPENSES_DB_DSN} \
		-smtp-host=${EXPENSES_SMTP_HOST} \
		-smtp-port=${EXPENSES_SMTP_PORT} \
		-smtp-username=${EXPENSES_SMTP_USERNAME} \
		-smtp-password=${EXPENSES_SMTP_PASSWORD} \
		-smtp-sender=${EXPENSES_SMTP_SENDER}

## db/migrations/new name=$1: create a new database migration
.PHONY: db/migrations/new
db/migrations/new:
	@echo "Creating migration files for ${name}..."
	migrate create -seq -ext=.sql -dir=./migrations ${name}
	
## db/migrations/up: apply all database migrations
.PHONY: db/migrations/up
db/migrations/up: confirm
	@echo "Running migrations..."
	migrate -path ./migrations -database="${EXPENSES-DB-DSN}" up

# ========================================= #
# Quality Control
# ========================================== #

## audit: tidy dependencies and format, vet and run unit tests
.PHONY: audit
audit:
	@echo "Tidying and veryfing module dependencies..."
	go mod tidy
	go mod verify
	@echo "Formatting code..."
	go fmt ./...
	@echo "Vetting code..."
	go vet ./...
	go test -v -short ./...
	
## test/integration: run integration tests
.PHONY: test/integration
test/integration:
	@echo "Starting service containers..."
	docker-compose -f docker-compose.test.yml up -d 
	@echo "Running tests..."
	go test ./internal/data
	@echo "Shutting down services..."

# ========================================= #
# Build
# ========================================== #

## build/api: build the cmd/api application
.PHONY: build/api
build/api:
	@echo "Building cmd/api..."
	go build -ldflags="-s" -o=./bin/api ./cmd/api
