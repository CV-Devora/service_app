-- +goose Up
-- +goose StatementBegin
ALTER TABLE pembelian
    ADD COLUMN IF NOT EXISTS berat_nota DECIMAL(10, 3) NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS harga_nota BIGINT NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS harga_rata BIGINT NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS tipe_pembayaran VARCHAR(50) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS jumlah_pembayaran BIGINT NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE pembelian
    DROP COLUMN IF EXISTS berat_nota,
    DROP COLUMN IF EXISTS harga_nota,
    DROP COLUMN IF EXISTS harga_rata,
    DROP COLUMN IF EXISTS tipe_pembayaran,
    DROP COLUMN IF EXISTS jumlah_pembayaran;
-- +goose StatementEnd
