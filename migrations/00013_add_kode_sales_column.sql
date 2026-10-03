-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
    ADD COLUMN IF NOT EXISTS kode_sales INT DEFAULT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Note: converting grup_id back to UUID is not reversible safely
-- +goose StatementEnd
