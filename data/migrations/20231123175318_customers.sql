-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE IF NOT EXISTS thyrasec.customers
(
    username character varying(255) COLLATE pg_catalog."default" NOT NULL,
    password_hash character varying(512) COLLATE pg_catalog."default" NOT NULL,
    email character varying(255) COLLATE pg_catalog."default",
    full_name character varying(255) COLLATE pg_catalog."default" NOT NULL,
    address character varying(512) COLLATE pg_catalog."default",
    phone_number character varying(50) COLLATE pg_catalog."default",
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    last_login timestamp without time zone,
    status user_status DEFAULT 'ACTIVE'::user_status,
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    customer_number character varying(100) COLLATE pg_catalog."default",
    CONSTRAINT customers_pkey PRIMARY KEY (id),
    CONSTRAINT customers_email_key UNIQUE (email),
    CONSTRAINT customers_username_key UNIQUE (username)
);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE thyrasec.customers
-- +goose StatementEnd
