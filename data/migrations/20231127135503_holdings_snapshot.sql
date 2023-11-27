-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS holdings_snapshots (
    snapshotID uuid NOT NULL DEFAULT uuid_generate_v4(),
    account_id uuid NOT NULL,
    asset_id uuid NOT NULL,
    quantity DECIMAL(15,4) NOT NULL,
    snapshot_date DATE,
    PRIMARY KEY (snapshotID),
    FOREIGN KEY (account_id) REFERENCES accounts(account_id),
    FOREIGN KEY (asset_id) REFERENCES assets(asset_id)
);

CREATE INDEX idx_snapshot_date ON holdings_snapshots(snapshot_date);
CREATE INDEX idx_account_id ON holdings_snapshots(account_id);
CREATE INDEX idx_asset_id ON holdings_snapshots(asset_id);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE holdings_snapshots;
-- +goose StatementEnd
