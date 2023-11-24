-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE IF NOT EXISTS thyrasec.asset_types
(
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    type_name character varying(100) COLLATE pg_catalog."default" NOT NULL,
    description text COLLATE pg_catalog."default",
    CONSTRAINT asset_types_pkey PRIMARY KEY (id),
    CONSTRAINT asset_types_type_name_key UNIQUE (type_name)
);
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE thyrasec.asset_types
-- +goose StatementEnd
