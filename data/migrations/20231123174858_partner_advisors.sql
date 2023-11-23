-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE IF NOT EXISTS thyrasec.partners_advisors
(
    username character varying(255) COLLATE pg_catalog."default" NOT NULL,
    password_hash character varying(512) COLLATE pg_catalog."default" NOT NULL,
    email character varying(255) COLLATE pg_catalog."default",
    full_name character varying(255) COLLATE pg_catalog."default" NOT NULL,
    company_name character varying(255) COLLATE pg_catalog."default",
    phone_number character varying(50) COLLATE pg_catalog."default",
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    last_login timestamp without time zone,
    status user_status DEFAULT 'ACTIVE'::user_status,
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    customer_number character varying(100) COLLATE pg_catalog."default",
    CONSTRAINT partners_advisors_pkey PRIMARY KEY (id),
    CONSTRAINT partners_advisors_email_key UNIQUE (email),
    CONSTRAINT partners_advisors_username_key UNIQUE (username)
);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE partners_advisors
-- +goose StatementEnd
