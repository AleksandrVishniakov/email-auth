-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    is_email_verified BOOLEAN NOT NULL DEFAULT false,
    email_verifying_hash VARCHAR(128) DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE If EXISTS users;
-- +goose StatementEnd
