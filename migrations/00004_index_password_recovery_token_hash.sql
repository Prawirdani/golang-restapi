-- +goose Up
-- +goose StatementBegin
SELECT
  'up SQL query';

CREATE INDEX idx_prt_token_hash ON password_recovery_tokens (token_hash);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
SELECT
  'down SQL query';

DROP INDEX IF EXISTS idx_prt_token_hash;

-- +goose StatementEnd
