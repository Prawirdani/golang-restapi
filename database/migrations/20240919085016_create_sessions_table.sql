-- +goose Up
-- +goose StatementBegin
SELECT
  'up SQL query';

CREATE TABLE IF NOT EXISTS sessions (
  id SERIAL PRIMARY KEY,
  user_id VARCHAR(100) NOT NULL,
  refresh_token VARCHAR(100) NOT NULL,
  user_agent VARCHAR(255) NOT NULL,
  expires_at TIMESTAMPTZ NOT NULL,
  accessed_at TIMESTAMPTZ NOT NULL,
  CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users (id)
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
SELECT
  'down SQL query';

DROP TABLE IF EXISTS sessions;

-- +goose StatementEnd
