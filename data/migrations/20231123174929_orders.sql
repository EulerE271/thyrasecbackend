-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE IF NOT EXISTS thyrasec.orders
(
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    account_id uuid NOT NULL,
    asset_id uuid NOT NULL,
    order_type text COLLATE pg_catalog."default" NOT NULL,
    quantity integer NOT NULL,
    price_per_unit double precision,
    total_amount double precision NOT NULL,
    status order_status_type NOT NULL DEFAULT 'created'::order_status_type,
    created_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT orders_pkey PRIMARY KEY (id),
    CONSTRAINT fk_account FOREIGN KEY (account_id)
        REFERENCES thyrasec.accounts (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE,
    CONSTRAINT fk_asset FOREIGN KEY (asset_id)
        REFERENCES thyrasec.assets (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
);
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE thyrasec.orders
-- +goose StatementEnd
