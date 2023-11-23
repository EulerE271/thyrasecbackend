-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE IF NOT EXISTS thyrasec.power_of_attorney
(
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    account_id integer,
    user_id integer,
    start_date date,
    end_date date,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    CONSTRAINT power_of_attorney_pkey PRIMARY KEY (id)
);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE power_of_attorney
-- +goose StatementEnd
