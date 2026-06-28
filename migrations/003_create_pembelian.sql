-- +goose Up
CREATE TABLE IF NOT EXISTS pembelian (
    id INTEGER PRIMARY KEY,
    no_faktur VARCHAR(255) NOT NULL UNIQUE,
    nama VARCHAR(255) NOT NULL,
    tipe_pemasok VARCHAR(100) NOT NULL,
    harga_deal INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP NULL
);

-- +goose Down
DROP TABLE IF EXISTS pembelian;
