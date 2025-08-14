-- +goose Up

CREATE TABLE IF NOT EXISTS roles (
    id uuid NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    role varchar(64) UNIQUE,
    created_at timestamptz DEFAULT NOW()
);

INSERT INTO roles(role) VALUES 
    ('admin'),
    ('user');

CREATE TABLE IF NOT EXISTS users_roles (
    id uuid NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    user_id uuid NOT NULL references users(id),
    role_id uuid NOT NULL references roles(id),
    created_at timestamptz DEFAULT NOW(),
    updated_at timestamptz DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_users_roles_user_id ON users_roles(user_id);

-- +goose Down

DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS users_roles;
