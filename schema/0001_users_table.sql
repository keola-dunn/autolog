-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id serial NOT NULL PRIMARY KEY,
    username text UNIQUE, 
    salt text, 
    password_hash text,
    email text UNIQUE,
    name text DEFAULT '',
    created_at timestamp DEFAULT NOW(),
    updated_at timestamp DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_users_username_email ON users(username, email);

-- +goose Down

DROP TABLE IF EXISTS users;
