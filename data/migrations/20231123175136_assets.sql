-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE IF NOT EXISTS thyrasec.assets
(
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    instrument_name character varying(255) COLLATE pg_catalog."default" NOT NULL,
    isin character varying(12) COLLATE pg_catalog."default" NOT NULL,
    ticker character varying(10) COLLATE pg_catalog."default" NOT NULL,
    exchange character varying(255) COLLATE pg_catalog."default",
    currency character varying(10) COLLATE pg_catalog."default",
    instrument_type character varying(100) COLLATE pg_catalog."default",
    current_price numeric(20,2),
    volume bigint,
    country character varying(100) COLLATE pg_catalog."default",
    sector character varying(100) COLLATE pg_catalog."default",
    asset_type_id uuid,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT assets_pkey PRIMARY KEY (id),
    CONSTRAINT assets_isin_key UNIQUE (isin),
    CONSTRAINT assets_asset_type_id_fkey FOREIGN KEY (asset_type_id)
        REFERENCES thyrasec.asset_types (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
);
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE thyrasec.assets
-- +goose StatementEnd
