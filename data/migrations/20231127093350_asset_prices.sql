-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE IF NOT EXISTS thyrasec.asset_prices
(
    asset_id uuid NOT NULL,
    price_date date NOT NULL,
    price numeric(20,2),
    CONSTRAINT fk_asset_price FOREIGN KEY (asset_id)
        REFERENCES thyrasec.assets (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE,
    PRIMARY KEY (asset_id, price_date)
);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
