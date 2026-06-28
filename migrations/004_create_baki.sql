-- +goose Up
CREATE TABLE IF NOT EXISTS baki (
    id UUID PRIMARY KEY,
    nama VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP NULL
);

-- +goose Down
DROP TABLE IF EXISTS baki;
