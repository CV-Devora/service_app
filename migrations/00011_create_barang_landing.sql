-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS barang_landing (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    nama        VARCHAR(255) NOT NULL,
    karat       VARCHAR(100),
    berat       DECIMAL(10, 3) DEFAULT 0,
    harga       BIGINT DEFAULT 0,
    photo       TEXT,
    created_at  TIMESTAMP DEFAULT NOW(),
    updated_at  TIMESTAMP DEFAULT NOW(),
    deleted_at  TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS barang_landing;
-- +goose StatementEnd
