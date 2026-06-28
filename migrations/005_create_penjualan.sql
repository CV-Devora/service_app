-- +goose Up
CREATE TABLE IF NOT EXISTS penjualan (
    id UUID PRIMARY KEY,
    no_faktur VARCHAR(255) NOT NULL UNIQUE,
    nama VARCHAR(255) NOT NULL,
    total_harga INT NOT NULL,
    kode_sales UUID NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP NULL,
    CONSTRAINT fk_penjualan_kode_sales
        FOREIGN KEY (kode_sales)
        REFERENCES users (id)
        ON UPDATE CASCADE
        ON DELETE SET NULL
);

-- +goose Down
DROP TABLE IF EXISTS penjualan;
