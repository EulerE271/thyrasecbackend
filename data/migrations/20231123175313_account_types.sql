-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE IF NOT EXISTS thyrasec.account_types
(
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    account_type_name character varying(255) COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT account_types_pkey PRIMARY KEY (id)
);
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE thyrasec.account_types
-- +goose StatementEnd
