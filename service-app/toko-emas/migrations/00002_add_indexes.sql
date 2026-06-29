-- +goose Up
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS idx_barang_deleted_at ON barang(deleted_at);
CREATE INDEX IF NOT EXISTS idx_barang_baki_id ON barang(baki_id);
CREATE INDEX IF NOT EXISTS idx_barang_pembelian_id ON barang(pembelian_id);
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);
CREATE INDEX IF NOT EXISTS idx_pembelian_deleted_at ON pembelian(deleted_at);
CREATE INDEX IF NOT EXISTS idx_karat_deleted_at ON karat(deleted_at);
CREATE INDEX IF NOT EXISTS idx_baki_deleted_at ON baki(deleted_at);
CREATE INDEX IF NOT EXISTS idx_penjualan_deleted_at ON penjualan(deleted_at);
CREATE INDEX IF NOT EXISTS idx_penjualan_kode_sales ON penjualan(kode_sales);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_barang_deleted_at;
DROP INDEX IF EXISTS idx_barang_baki_id;
DROP INDEX IF EXISTS idx_barang_pembelian_id;
DROP INDEX IF EXISTS idx_users_deleted_at;
DROP INDEX IF EXISTS idx_pembelian_deleted_at;
DROP INDEX IF EXISTS idx_karat_deleted_at;
DROP INDEX IF EXISTS idx_baki_deleted_at;
DROP INDEX IF EXISTS idx_penjualan_deleted_at;
DROP INDEX IF EXISTS idx_penjualan_kode_sales;
-- +goose StatementEnd
