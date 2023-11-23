-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE IF NOT EXISTS thyrasec.transactions
(
    id uuid NOT NULL,
    type uuid NOT NULL,
    asset1_id uuid,
    asset2_id uuid,
    amount_asset1 integer,
    amount_asset2 integer,
    created_by_id uuid NOT NULL,
    updated_by_id uuid NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    corrected boolean NOT NULL,
    canceled boolean NOT NULL,
    status_transaction uuid NOT NULL,
    comment character varying(100) COLLATE pg_catalog."default",
    transaction_owner_id uuid,
    account_owner_id uuid,
    account_asset1_id uuid,
    account_asset2_id uuid,
    trade_date date,
    settlement_date date,
    asset2_currency uuid,
    asset2_price double precision,
    order_no text COLLATE pg_catalog."default",
    business_event uuid,
    CONSTRAINT transactions_pkey PRIMARY KEY (id)
);
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
DROP TABLE thyrasec.transactions