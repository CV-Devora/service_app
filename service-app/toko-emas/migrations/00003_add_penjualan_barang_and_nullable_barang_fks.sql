-- +goose Up
-- +goose StatementBegin
ALTER TABLE barang
    ALTER COLUMN pembelian_id DROP NOT NULL,
    ALTER COLUMN baki_id DROP NOT NULL;

ALTER TABLE penjualan
    ALTER COLUMN kode_sales DROP NOT NULL;

CREATE TABLE IF NOT EXISTS penjualan_barang (
    barang_id UUID NOT NULL REFERENCES barang(id) ON DELETE CASCADE,
    penjualan_id UUID NOT NULL REFERENCES penjualan(id) ON DELETE CASCADE,
    PRIMARY KEY (barang_id, penjualan_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS penjualan_barang;

ALTER TABLE penjualan
    ALTER COLUMN kode_sales SET NOT NULL;

ALTER TABLE barang
    ALTER COLUMN pembelian_id SET NOT NULL,
    ALTER COLUMN baki_id SET NOT NULL;
-- +goose StatementEnd
