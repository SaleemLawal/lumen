MIGRATIONS_PATH = ./backend/cmd/migrate/migrations
DB_URL = postgres://admin:password@localhost:5432/lumen?sslmode=disable
.PHONY: swagger
swagger:
	swag init -g api.go -d ./backend/cmd/api -o ./backend/docs --parseInternal --parseDependency

.PHONY: migration
migration:
	@echo "Creating migration..."
	migrate create -seq -ext sql -dir ${MIGRATIONS_PATH} $(filter-out $@,$(MAKECMDGOALS))

.PHONY: migrate-up
migrate-up:
	@echo "Migrating up..."
	migrate -path ${MIGRATIONS_PATH} -database ${DB_URL} up $(filter-out $@,$(MAKECMDGOALS))

.PHONY: migrate-down
migrate-down:
	@echo "Migrating down..."
	migrate -path ${MIGRATIONS_PATH} -database ${DB_URL} down $(filter-out $@,$(MAKECMDGOALS))

.PHONY: migrate-force
migrate-force:
	@echo "Migrating force..."
	migrate -path ${MIGRATIONS_PATH} -database ${DB_URL} force $(filter-out $@,$(MAKECMDGOALS))
	@make migrate-up