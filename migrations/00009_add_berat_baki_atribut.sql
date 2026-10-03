-- +goose Up
-- +goose StatementBegin
ALTER TABLE baki
    ADD COLUMN IF NOT EXISTS berat DECIMAL(10, 3) NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- +goose StatementEnd
