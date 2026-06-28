-- +goose Up
CREATE TABLE IF NOT EXISTS barang (
    id UUID PRIMARY KEY,
    barcode VARCHAR(255) UNIQUE,
    nama VARCHAR(255) NOT NULL,
    karat INT NOT NULL,
    berat DECIMAL(18, 2) NOT NULL,
    harga INT NOT NULL,
    photo VARCHAR(255),
    kondisi VARCHAR(100) NOT NULL,
    pembelian_id INTEGER NULL,
    baki_id UUID NULL,
    grup_id UUID NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP NULL,
    CONSTRAINT fk_barang_pembelian
        FOREIGN KEY (pembelian_id)
        REFERENCES pembelian (id)
        ON UPDATE CASCADE
        ON DELETE SET NULL,
    CONSTRAINT fk_barang_baki
        FOREIGN KEY (baki_id)
        REFERENCES baki (id)
        ON UPDATE CASCADE
        ON DELETE SET NULL
);

-- +goose Down
DROP TABLE IF EXISTS barang;
