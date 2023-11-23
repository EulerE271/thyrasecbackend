-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE IF NOT EXISTS thyrasec.holdings
(
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    account_id uuid NOT NULL,
    asset_id uuid NOT NULL,
    quantity integer NOT NULL,
    CONSTRAINT holdings_pkey PRIMARY KEY (id),
    CONSTRAINT fk_account FOREIGN KEY (account_id)
        REFERENCES thyrasec.accounts (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE,
    CONSTRAINT fk_asset FOREIGN KEY (asset_id)
        REFERENCES thyrasec.assets (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE,
    CONSTRAINT holdings_quantity_check CHECK (quantity >= 0)
);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE thyrasec.holdings
-- +goose StatementEnd
