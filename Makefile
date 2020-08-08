include .env
export

.PHONY: migrate-new
migrate-new: ## create a new database migration
	@read -p "Enter the name of the new migration: " name; \
	sql-migrate new $$name;

.PHONY: migrate-up
migrate-up:
	@echo "Running database migration..."
	@sql-migrate up

.PHONY: migrate-down
migrate-down:
	@echo "Undoing last applied migration..."
	@sql-migrate down

.PHONY: migrate-fresh
migrate-fresh:
	@echo "Resetting database..."
	@sudo mysql -u"$$DB_USERNAME" -p"$$DB_PASSWORD" -e "DROP DATABASE IF EXISTS $$DB_NAME; CREATE DATABASE $$DB_NAME"
	@echo "Running migration..."
	@sql-migrate up
	
test-all:
	@go test -v ./... -tags=all