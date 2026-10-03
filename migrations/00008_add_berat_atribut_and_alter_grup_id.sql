-- +goose Up
-- +goose StatementBegin
ALTER TABLE barang
    ADD COLUMN IF NOT EXISTS berat_atribut DECIMAL(10, 3) NOT NULL DEFAULT 0;

ALTER TABLE barang
    ALTER COLUMN grup_id TYPE VARCHAR(255) USING COALESCE(grup_id::TEXT, ''),
    ALTER COLUMN grup_id SET DEFAULT '';

ALTER TABLE barang
    RENAME COLUMN grup_id TO grup;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE barang DROP COLUMN IF EXISTS berat_atribut;
-- Note: converting grup_id back to UUID is not reversible safely
-- +goose StatementEnd
