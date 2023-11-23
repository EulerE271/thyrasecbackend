-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE IF NOT EXISTS thyrasec.accounts
(
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    account_name character varying(100) COLLATE pg_catalog."default" NOT NULL,
    account_type uuid NOT NULL,
    account_owner_company boolean NOT NULL,
    account_balance double precision NOT NULL DEFAULT 0,
    account_currency character varying(3) COLLATE pg_catalog."default" NOT NULL,
    account_number character varying(20) COLLATE pg_catalog."default",
    account_status character varying(20) COLLATE pg_catalog."default",
    interest_rate double precision,
    overdraft_limit double precision,
    account_description text COLLATE pg_catalog."default",
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    created_by character varying(100) COLLATE pg_catalog."default",
    updated_by character varying(100) COLLATE pg_catalog."default",
    account_holder_id uuid NOT NULL,
    available_cash numeric(10,2) NOT NULL DEFAULT 0,
    reserved_cash numeric(10,2) NOT NULL DEFAULT 0,
    CONSTRAINT accounts_pkey PRIMARY KEY (id)
);
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE thyrasec.accounts
-- +goose StatementEnd
