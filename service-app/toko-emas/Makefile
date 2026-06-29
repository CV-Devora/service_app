.PHONY: run migrate migrate-down migrate-status swagger tidy

# Jalankan server
run:
	go run main.go -conf configs

# Jalankan migrasi database (up)
migrate:
	go run cmd/migrate/main.go -conf configs -cmd up

# Rollback migrasi
migrate-down:
	go run cmd/migrate/main.go -conf configs -cmd down

# Status migrasi
migrate-status:
	go run cmd/migrate/main.go -conf configs -cmd status

# Reset semua migrasi
migrate-reset:
	go run cmd/migrate/main.go -conf configs -cmd reset

# Generate swagger docs (perlu install: go install github.com/swaggo/swag/cmd/swag@latest)
swagger:
	swag init -g main.go --output docs

# Download dependencies
tidy:
	go mod tidy

# Build binary
build:
	go build -o bin/server main.go
	go build -o bin/migrate cmd/migrate/main.go
