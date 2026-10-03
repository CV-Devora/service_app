# Toko Emas API

REST API untuk manajemen toko emas dibangun dengan **Go-Kratos** pattern, **GORM + PostgreSQL**, **Goose** migration, dan **Swagger** documentation.

## Tech Stack

- **Framework**: Go + Gorilla Mux (HTTP layer ala Kratos)
- **Database**: PostgreSQL + GORM
- **Migration**: Goose
- **Docs**: Swagger UI (`/docs/index.html`)
- **Config**: YAML (`configs/config.yaml`)

## Struktur Proyek

```
toko-emas/
├── main.go                     # Entry point (go run main.go -conf configs)
├── Makefile
├── configs/
│   └── config.yaml             # Konfigurasi DB & server
├── migrations/
│   ├── 00001_create_tables.sql
│   └── 00002_add_indexes.sql
├── api/v1/
│   └── types.go                # Request/Response types
├── docs/
│   └── docs.go                 # Swagger generated docs
├── internal/
│   ├── conf/conf.go            # Config structs
│   ├── data/
│   │   ├── db.go               # Database connection
│   │   ├── model.go            # GORM models
│   │   ├── barang_repo.go
│   │   ├── user_repo.go
│   │   ├── pembelian_repo.go
│   │   └── other_repos.go      # Karat, Baki, Penjualan repos
│   ├── service/
│   │   ├── barang_service.go
│   │   ├── user_service.go
│   │   └── other_services.go
│   └── server/
│       └── http.go             # Router setup
└── cmd/
    ├── server/main.go
    └── migrate/main.go         # Migrate CLI
```

## Setup & Menjalankan

### 1. Prasyarat

- Go 1.22+
- PostgreSQL
- (Opsional) `swag` untuk regenerate docs

### 2. Install Dependencies

```bash
go mod tidy
```

### 3. Konfigurasi Database

Edit `configs/config.yaml`:

```yaml
data:
  database:
    dsn: "host=localhost user=postgres password=postgres dbname=toko_emas port=5432 sslmode=disable"
```

### 4. Buat Database

```sql
CREATE DATABASE toko_emas;
```

### 5. Jalankan Migrasi

```bash
# Migrasi semua tabel
make migrate
# atau manual:
go run cmd/migrate/main.go -conf configs -cmd up

# Cek status migrasi
make migrate-status

# Rollback 1 step
make migrate-down
```

### 6. Jalankan Server

```bash
# Cara utama
go run main.go -conf configs

# Atau via Makefile
make run
```

Server berjalan di `http://localhost:8000`

## API Endpoints

| Method | Endpoint                 | Deskripsi              |
|--------|--------------------------|------------------------|
| GET    | /health                  | Health check           |
| GET    | /docs/index.html         | Swagger UI             |
| GET    | /api/v1/barang           | List barang            |
| POST   | /api/v1/barang           | Tambah barang          |
| GET    | /api/v1/barang/{id}      | Detail barang          |
| PUT    | /api/v1/barang/{id}      | Update barang          |
| DELETE | /api/v1/barang/{id}      | Hapus barang           |
| GET    | /api/v1/users            | List users             |
| POST   | /api/v1/users            | Tambah user            |
| GET    | /api/v1/users/{id}       | Detail user            |
| PUT    | /api/v1/users/{id}       | Update user            |
| DELETE | /api/v1/users/{id}       | Hapus user             |
| GET    | /api/v1/pembelian        | List pembelian         |
| POST   | /api/v1/pembelian        | Tambah pembelian       |
| GET    | /api/v1/pembelian/{id}   | Detail pembelian       |
| PUT    | /api/v1/pembelian/{id}   | Update pembelian       |
| DELETE | /api/v1/pembelian/{id}   | Hapus pembelian        |
| GET    | /api/v1/karat            | List karat             |
| POST   | /api/v1/karat            | Tambah karat           |
| GET    | /api/v1/karat/{id}       | Detail karat           |
| PUT    | /api/v1/karat/{id}       | Update karat           |
| DELETE | /api/v1/karat/{id}       | Hapus karat            |
| GET    | /api/v1/baki             | List baki              |
| POST   | /api/v1/baki             | Tambah baki            |
| GET    | /api/v1/baki/{id}        | Detail baki            |
| PUT    | /api/v1/baki/{id}        | Update baki            |
| DELETE | /api/v1/baki/{id}        | Hapus baki             |
| GET    | /api/v1/penjualan        | List penjualan         |
| POST   | /api/v1/penjualan        | Tambah penjualan       |
| GET    | /api/v1/penjualan/{id}   | Detail penjualan       |
| PUT    | /api/v1/penjualan/{id}   | Update penjualan       |
| DELETE | /api/v1/penjualan/{id}   | Hapus penjualan        |

## Contoh Request

### Tambah Barang

```bash
curl -X POST http://localhost:8000/api/v1/barang \
  -H "Content-Type: application/json" \
  -d '{
    "barcode": "BC001",
    "nama": "Gelang Emas 24K",
    "karat": 24,
    "berat": 5.5,
    "harga": 5500000,
    "kondisi": "baru"
  }'
```

### Tambah User

```bash
curl -X POST http://localhost:8000/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "nama": "Budi Santoso",
    "username": "budi",
    "password": "secret123",
    "role": "kasir"
  }'
```

## Regenerate Swagger Docs

```bash
# Install swag
go install github.com/swaggo/swag/cmd/swag@latest

# Generate
make swagger
# atau
swag init -g main.go --output docs
```

## Migration Commands

```bash
go run cmd/migrate/main.go -conf configs -cmd up      # Jalankan semua migrasi
go run cmd/migrate/main.go -conf configs -cmd down    # Rollback 1 step
go run cmd/migrate/main.go -conf configs -cmd status  # Lihat status
go run cmd/migrate/main.go -conf configs -cmd reset   # Reset semua
```
