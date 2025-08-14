-- +goose Up

CREATE TABLE IF NOT EXISTS shops (
    id uuid NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    name varchar(256),
    address1 varchar(256),
    address2 varchar(256),
    city varchar(128),
    state varchar(4),
    zip varchar(12),
    phone varchar(20),
    
    created_by uuid NOT NULL references users(id),
    created_at timestamptz DEFAULT NOW(),
    updated_at timestamptz DEFAULT NOW()
);

INSERT INTO roles(role) VALUES 
    ('mechanic');

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
CREATE INDEX IF NOT EXISTS idx_cars_vin ON cars(vin);

CREATE TABLE IF NOT EXISTS users_cars (
    id uuid NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    user_id uuid NOT NULL references users(id),
    car_id uuid NOT NULL references cars(id),
    created_at timestamptz DEFAULT NOW(),
    updated_at timestamptz DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_users_cars_user_id ON users_cars(user_id);


CREATE TABLE IF NOT EXISTS service_logs (
    id uuid NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    user_id uuid NOT NULL references users(id),
    car_id uuid NOT NULL references cars(id),

    "type" varchar(128),
    "date" date,
    mileage integer,
    details JSONB,
    notes text,

    created_at timestamptz DEFAULT NOW(),
    updated_at timestamptz DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_service_logs_user_id ON service_logs(user_id);

-- +goose Down

DROP TABLE IF EXISTS service_logs;
DROP TABLE IF EXISTS users_cars;
DROP TABLE IF EXISTS cars;
DROP TABLE IF EXISTS shops;
