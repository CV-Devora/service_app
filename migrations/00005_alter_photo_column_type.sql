-- +goose Up
-- +goose StatementBegin
ALTER TABLE barang ALTER COLUMN photo TYPE TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE barang ALTER COLUMN photo TYPE VARCHAR(500);
-- +goose StatementEnd
