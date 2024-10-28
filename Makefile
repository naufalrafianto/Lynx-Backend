# Makefile
.PHONY: dev prod test clean migrate-up migrate-down create-migration

# Development
dev:
	docker-compose -f docker-compose.dev.yml up --build

# Production
prod:
	docker-compose -f docker-compose.prod.yml up --build -d

# Testing
test:
	go test ./... -v

# Cleanup
clean:
	docker-compose -f docker-compose.dev.yml down -v
	docker-compose -f docker-compose.prod.yml down -v
	rm -rf tmp/

# Linting
lint:
	golangci-lint run

# Dependencies
tidy:
	go mod tidy

# Migration commands
migrate-up:
	@echo "Running migrations up..."
	@docker-compose -f docker-compose.dev.yml exec -T app go run cmd/migrate/main.go -command=up

migrate-down:
	@echo "Running migrations down..."
	@docker-compose -f docker-compose.dev.yml exec -T app go run cmd/migrate/main.go -command=down

# Create new migration
create-migration:
	@read -p "Enter migration name: " name; \
	timestamp=$$(date +%Y%m%d%H%M%S); \
	mkdir -p migrations; \
	touch migrations/$${timestamp}_$$name.up.sql; \
	touch migrations/$${timestamp}_$$name.down.sql; \
	echo "Created migrations/$${timestamp}_$$name.up.sql"; \
	echo "Created migrations/$${timestamp}_$$name.down.sql"