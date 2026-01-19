.PHONY: run test migrate up down

up:
	docker-compose up -d --build

down:
	docker-compose down

run:
	go run cmd/api/main.go

test:
	go test -v ./...

migrate:
	# This is a placeholder. In a real scenario, you'd use a tool like migrate/migrate
	# For simplicity in this tech task, we might just load the SQL file on app startup or use a simple script.
	# But let's assume we have a migration tool or we'll run it via psql for local dev.
	@echo "Running migrations..."
