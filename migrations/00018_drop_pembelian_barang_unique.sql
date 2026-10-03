-- +goose Up
-- +goose StatementBegin
ALTER TABLE barang
    DROP CONSTRAINT IF EXISTS uq_barang_pembelian_id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE barang
    ADD CONSTRAINT uq_barang_pembelian_id UNIQUE (pembelian_id);
-- +goose StatementEnd