-- +goose Up
-- +goose StatementBegin
ALTER TABLE penjualan
    ADD COLUMN IF NOT EXISTS harga_gram BIGINT NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS harga_jual BIGINT NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS ongkos BIGINT NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS harga_total BIGINT NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS cash BIGINT NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS transfer BIGINT NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS debet BIGINT NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Note: converting grup_id back to UUID is not reversible safely
-- +goose StatementEnd
