-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    nama        VARCHAR(255),
    username    VARCHAR(255) UNIQUE,
    password    VARCHAR(255),
    role        VARCHAR(50),
    created_at  TIMESTAMP DEFAULT NOW(),
    updated_at  TIMESTAMP DEFAULT NOW(),
    deleted_at  TIMESTAMP
);

CREATE TABLE IF NOT EXISTS pembelian (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    no_faktur       VARCHAR(100) UNIQUE,
    nama            VARCHAR(255),
    tipe_pemasok    VARCHAR(100),
    harga_deal      BIGINT,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW(),
    deleted_at      TIMESTAMP
);

CREATE TABLE IF NOT EXISTS karat (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name        VARCHAR(100),
    harga       BIGINT,
    created_at  TIMESTAMP DEFAULT NOW(),
    updated_at  TIMESTAMP DEFAULT NOW(),
    deleted_at  TIMESTAMP
);

CREATE TABLE IF NOT EXISTS baki (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    nama        VARCHAR(255),
    created_at  TIMESTAMP DEFAULT NOW(),
    updated_at  TIMESTAMP DEFAULT NOW(),
    deleted_at  TIMESTAMP
);

CREATE TABLE IF NOT EXISTS barang (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    barcode         VARCHAR(100) UNIQUE,
    nama            VARCHAR(255),
    karat           INT,
    berat           DECIMAL(10, 3),
    harga           BIGINT,
    photo           VARCHAR(500),
    kondisi         VARCHAR(100),
    pembelian_id    UUID REFERENCES pembelian(id) ON DELETE SET NULL,
    baki_id         UUID REFERENCES baki(id) ON DELETE SET NULL,
    grup_id         UUID,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW(),
    deleted_at      TIMESTAMP
);

CREATE TABLE IF NOT EXISTS penjualan (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    no_faktur   VARCHAR(100) UNIQUE,
    nama        VARCHAR(255),
    total_harga BIGINT,
    kode_sales  UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at  TIMESTAMP DEFAULT NOW(),
    updated_at  TIMESTAMP DEFAULT NOW(),
    deleted_at  TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS penjualan;
DROP TABLE IF EXISTS barang;
DROP TABLE IF EXISTS baki;
DROP TABLE IF EXISTS karat;
DROP TABLE IF EXISTS pembelian;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
