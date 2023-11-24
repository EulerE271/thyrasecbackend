-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE IF NOT EXISTS thyrasec.cash_reservations
(
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    order_id uuid NOT NULL,
    account_id uuid NOT NULL,
    amount numeric(15,2) NOT NULL,
    reserved_until timestamp without time zone NOT NULL,
    status character varying(255) COLLATE pg_catalog."default" NOT NULL,
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    updated_at timestamp without time zone NOT NULL DEFAULT now(),
    CONSTRAINT cash_reservations_pkey PRIMARY KEY (id),
    CONSTRAINT fk_account FOREIGN KEY (account_id)
        REFERENCES thyrasec.accounts (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE,
    CONSTRAINT fk_order FOREIGN KEY (order_id)
        REFERENCES thyrasec.orders (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE thyrasec.cash_reservations
-- +goose StatementEnd
