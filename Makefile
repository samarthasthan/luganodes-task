
up:
	@echo "Running in Development mode..."
	@docker compose -f ./build/compose/compose.yaml up -d
	@echo "Development mode completed."

	
down:
	@echo "Running in Development mode..."
	@docker compose -f ./build/compose/compose.yaml down --volumes
	@echo "Development mode completed."


# Make migrations
migrate-up:
	@echo "Making migrations..."
	@migrate -path ./internal/store/database/mysql/migrations -database "mysql://root:${MYSQL_ROOT_PASSWORD}@tcp(${MYSQL_HOST}:${MYSQL_PORT})/luganodes" -verbose up
	@echo "Migrations completed."

# Delete migrations
migrate-down:
	@echo "Deleting migrations..."
	@migrate -path ./internal/store/database/mysql/migrations -database "mysql://root:${MYSQL_ROOT_PASSWORD}@tcp(${MYSQL_HOST}:${MYSQL_PORT})/luganodes" -verbose down

# Make migrations Devlopment environment
migrate-up-dev:
	@echo "Making migrations..."
	@migrate -path ./internal/store/database/mysql/migrations -database "mysql://root:password@tcp(localhost:3306)/luganodes" -verbose up
	@echo "Migrations completed."

# Delete migrations Devlopment environment
migrate-down-dev:
	@echo "Deleting migrations..."
	@migrate -path ./internal/store/database/mysql/migrations -database "mysql://root:password@tcp(localhost:3306)/luganodes" -verbose down

# SQLC generate
sqlc-gen:
	@echo "Generating SQLC..."
	@sqlc generate -f ./internal/store/database/mysql/sqlc/sqlc.yaml
	@echo "SQLC generation completed."