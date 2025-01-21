-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id serial NOT NULL PRIMARY KEY,
    username text,
    salt text, 
    salty_password text,
    email text,
    name text,
    created_at timestamp DEFAULT NOW(),
    updated_at timestamp DEFAULT NOW()
);

-- +goose Down
DROP TABLE users;