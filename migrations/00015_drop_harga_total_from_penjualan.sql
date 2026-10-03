-- +goose Up
-- +goose StatementBegin
ALTER TABLE penjualan
    DROP COLUMN IF EXISTS harga_total;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE penjualan
    ADD COLUMN IF NOT EXISTS harga_total BIGINT NOT NULL DEFAULT 0;
-- +goose StatementEnd
