-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE IF NOT EXISTS thyrasec.account_snapshots
(
    account_id uuid NOT NULL,
    snapshot_date date NOT NULL,
    total_value numeric(20,2),
    CONSTRAINT fk_account_snapshot FOREIGN KEY (account_id)
        REFERENCES thyrasec.accounts (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE,
    PRIMARY KEY (account_id, snapshot_date)
);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
