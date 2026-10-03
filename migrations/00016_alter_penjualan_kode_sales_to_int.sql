-- +goose Up
-- +goose StatementBegin
ALTER TABLE penjualan
    DROP CONSTRAINT IF EXISTS penjualan_kode_sales_fkey;

ALTER TABLE penjualan
    ADD COLUMN kode_sales_new INT;

UPDATE penjualan p
    SET kode_sales_new = u.kode_sales
    FROM users u
    WHERE u.id = p.kode_sales;

ALTER TABLE penjualan
    DROP COLUMN IF EXISTS kode_sales;

ALTER TABLE penjualan
    RENAME COLUMN kode_sales_new TO kode_sales;

CREATE INDEX IF NOT EXISTS idx_penjualan_kode_sales ON penjualan(kode_sales);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE penjualan
    ADD COLUMN kode_sales_old UUID REFERENCES users(id) ON DELETE SET NULL;

UPDATE penjualan p
    SET kode_sales_old = u.id
    FROM users u
    WHERE u.kode_sales = p.kode_sales;

ALTER TABLE penjualan
    DROP COLUMN IF EXISTS kode_sales;

ALTER TABLE penjualan
    RENAME COLUMN kode_sales_old TO kode_sales;

CREATE INDEX IF NOT EXISTS idx_penjualan_kode_sales ON penjualan(kode_sales);
-- +goose StatementEnd