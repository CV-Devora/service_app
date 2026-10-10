-- +goose Up
-- +goose StatementBegin
INSERT INTO users (id, nama, username, password, role, created_at, updated_at)
VALUES 
    (uuid_generate_v4(), 'Administrator', 'admin', '$2a$10$gVaSI4kaArrGOkt0bcTtA.rBpfgxo4Gnbqimb6nCmW0vB9gLVeR36', 'admin', NOW(), NOW())
ON CONFLICT (username) DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM users WHERE username IN ('admin', 'kasir1', 'sales1', 'manager1', 'owner1');
-- +goose StatementEnd
