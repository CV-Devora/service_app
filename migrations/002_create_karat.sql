-- +goose Up
CREATE TABLE IF NOT EXISTS karat (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    harga INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP NULL
);

-- +goose Down
DROP TABLE IF EXISTS karat;
