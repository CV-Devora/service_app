-- +goose Up
-- +goose StatementBegin

-- Hapus tabel weak entity penjualan_barang
DROP TABLE IF EXISTS penjualan_barang;

-- Tambahkan kolom penjualan_id pada barang (one-to-one dengan penjualan)
ALTER TABLE barang
    ADD COLUMN penjualan_id UUID REFERENCES penjualan(id) ON DELETE SET NULL;

-- UNIQUE constraint untuk menjamin relasi one-to-one
ALTER TABLE barang
    ADD CONSTRAINT uq_barang_penjualan_id UNIQUE (penjualan_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Hapus kolom penjualan_id dari barang
ALTER TABLE barang
    DROP CONSTRAINT IF EXISTS uq_barang_penjualan_id,
    DROP COLUMN IF EXISTS penjualan_id;

-- Buat ulang tabel penjualan_barang
CREATE TABLE IF NOT EXISTS penjualan_barang (
    barang_id    UUID NOT NULL REFERENCES barang(id) ON DELETE CASCADE,
    penjualan_id UUID NOT NULL REFERENCES penjualan(id) ON DELETE CASCADE,
    PRIMARY KEY (barang_id, penjualan_id)
);

-- +goose StatementEnd
