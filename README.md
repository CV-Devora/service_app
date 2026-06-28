# service-app

Scaffold backend Go-Kratos dengan clean architecture, PostgreSQL, `sqlc`, JSON API, dan Swagger/OpenAPI.

## Struktur

- `cmd/api`: entrypoint service
- `internal/config`: konfigurasi environment
- `internal/domain`: entity dan contract
- `internal/usecase`: business logic
- `internal/infrastructure/postgres`: repository PostgreSQL
- `internal/transport/http`: handler dan routing
- `db/schema.sql`: skema database
- `db/queries`: query `sqlc`
- `api/openapi.yaml`: dokumentasi API

## Setup

1. Siapkan PostgreSQL dan set `DATABASE_URL`.
2. Generate kode `sqlc`:
   `sqlc generate`
3. Jalankan service:
   `go run ./cmd/api`

Saat service start, migrasi PostgreSQL dijalankan otomatis lewat `goose`.

## Endpoint

- `GET /healthz`
- `GET /users`
- `POST /users`
- `GET /users/{id}`
- `PUT /users/{id}`
- `DELETE /users/{id}`

## Swagger

OpenAPI tersedia di:

- `GET /swagger/openapi.yaml`
