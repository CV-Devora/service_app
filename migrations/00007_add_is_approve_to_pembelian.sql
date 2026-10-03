-- +goose Up
-- +goose StatementBegin
ALTER TABLE pembelian ADD COLUMN IF NOT EXISTS is_approve BOOLEAN DEFAULT FALSE;
ALTER TABLE penjualan DROP COLUMN IF EXISTS is_approve;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE pembelian DROP COLUMN IF EXISTS is_approve;
ALTER TABLE penjualan ADD COLUMN IF NOT EXISTS is_approve BOOLEAN DEFAULT FALSE;
-- +goose StatementEnd
