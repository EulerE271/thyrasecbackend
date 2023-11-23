-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE IF NOT EXISTS thyrasec.transaction_logs
(
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    status character varying(255) COLLATE pg_catalog."default" NOT NULL,
    message text COLLATE pg_catalog."default",
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    transaction_id1 uuid NOT NULL,
    transaction_id2 uuid NOT NULL,
    CONSTRAINT transaction_logs_pkey PRIMARY KEY (id)
);
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE thyrasec.transaction_logs
-- +goose StatementEnd
