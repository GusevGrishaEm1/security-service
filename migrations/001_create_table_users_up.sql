-- +goose Up
CREATE TABLE IF NOT EXISTS "users" (
    id INTEGER PRIMARY KEY,
    email TEXT NOT NULL,
    "password" TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS "users";