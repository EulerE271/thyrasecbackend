-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TYPE user_status AS ENUM ('ACTIVE', 'INACTIVE', 'SUSPENDED', 'DELETED');

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TYPE IF EXISTS user_status;
-- +goose StatementEnd
