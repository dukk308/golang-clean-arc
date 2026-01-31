# Load environment variables from .env file if it exists
ifneq (,$(wildcard .env))
include .env
export $(shell sed 's/=.*//' .env)
else
$(info Warning: .env file not found. Please create .env file with your environment variables.)
$(info You can run 'go run main.go outenv' to see all required environment variables.)
endif
MIGRATION_DIR=database/migrations

migration-create-%:
	@goose -dir $(MIGRATION_DIR) create $* sql
migration-up:
	@goose -dir $(MIGRATION_DIR) postgres "$(DB_DSN)" up
migrate-down:
	@goose -dir $(MIGRATION_DIR) postgres "$(DB_DSN)" down
migration-gen:
	@echo "❌ Missing migration name. Usage: make migration-gen-<name>"
	@exit 1
migration-gen-%:
	@echo "🔧 Generating migration: $*"
	@atlas migrate diff $* --env gorm
migration-hash:
	@atlas migrate hash --env gorm
migration-status:
	@echo "🔍 Checking migration status with Atlas"
	@atlas migrate status --env gorm
migration-validate:
	@atlas migrate validate --env gorm

migrate-hash:
	@atlas migrate hash --env gorm
migrate-down-1:
	@goose -dir $(MIGRATION_DIR) postgres "$(DB_DSN)" down 1

swag:
	@swag init --parseDependency --output ./api-docs/swagger --parseInternal
