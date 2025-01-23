-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id serial NOT NULL PRIMARY KEY,
    username text UNIQUE,
    salt text, 
    password_hash text,
    email text,
    name text DEFAULT '',
    created_at timestamp DEFAULT NOW(),
    updated_at timestamp DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- +goose Down

DROP TABLE IF EXISTS users;
