-- +goose Up

INSERT INTO roles(role) VALUES 
    ('mechanic');

CREATE TABLE IF NOT EXISTS employees (
    id uuid NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    user_id uuid NOT NULL references users(id),
    shop_id uuid NOT NULL references shops(id),
    "role" varchar(64),
    created_by uuid NOT NULL references users(id),
    created_at timestamptz DEFAULT NOW(),
    updated_at timestamptz DEFAULT NOW()
);

-- +goose Down

DROP TABLE IF EXISTS employees;
DELETE FROM roles WHERE role = 'mechanic';
