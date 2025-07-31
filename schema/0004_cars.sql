-- +goose Up

CREATE TABLE IF NOT EXISTS cars (
    id uuid NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,

    make varchar(256),
    model varchar(256),
    trim varchar(256),
    year smallint,
    vin varchar(64),

    created_at timestamptz DEFAULT NOW(),
    updated_at timestamptz DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS users_cars (
    id uuid NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    user_id uuid NOT NULL references users(id),
    car_id uuid NOT NULL references cars(id),
    created_at timestamptz DEFAULT NOW(),
    updated_at timestamptz DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_users_cars_user_id ON users_cars(user_id);

-- +goose Down

DROP TABLE IF EXISTS users_cars;
DROP TABLE IF EXISTS cars;
