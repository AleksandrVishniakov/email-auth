-- +goose Up
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN email_verifying_hash;
ALTER TABLE users ADD COLUMN email_verifying_code INTEGER;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN email_verifying_hash VARCHAR(128);
ALTER TABLE users DROP COLUMN email_verifying_code;
-- +goose StatementEnd
