-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE IF NOT EXISTS thyrasec.transactions_types
(
    type_id uuid NOT NULL DEFAULT uuid_generate_v4(),
    transaction_type_name character varying(255) COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT transactions_types_pkey PRIMARY KEY (type_id)
);
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
DROP TABLE thyrasec.transactions