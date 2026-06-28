DB_URL ?= postgres://postgres:postgres@localhost:5432/jason_jewelry?sslmode=disable
PSH := pwsh

.PHONY: migrate-up migrate-down

migrate-up:
	$(PSH) -NoProfile -ExecutionPolicy Bypass -File scripts/migrate.ps1 -Action up -DatabaseUrl "$(DB_URL)"

migrate-down:
	$(PSH) -NoProfile -ExecutionPolicy Bypass -File scripts/migrate.ps1 -Action down -DatabaseUrl "$(DB_URL)"
