-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id uuid NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    simple_id serial NOT NULL,
    username varchar(128) UNIQUE, 
    salt text, 
    password_hash text,
    email varchar(256) UNIQUE, --email addresses should be limited to 254 or so
    name text DEFAULT '',
    created_at timestamptz DEFAULT NOW(),
    updated_at timestamptz DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_users_username_email ON users(username, email);

-- +goose Down

DROP TABLE IF EXISTS users;
