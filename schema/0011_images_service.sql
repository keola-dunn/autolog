-- +goose Up
CREATE SCHEMA IF NOT EXISTS images;

CREATE TABLE IF NOT EXISTS images.images (
    id uuid NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    user_id uuid NOT NULL,
    title varchar(256),
    path varchar(256) NOT NULL,
    width integer, 
    height integer,
    imageSizeKb integer, 
    hash text,
    created_at timestamptz DEFAULT NOW(),
    updated_at timestamptz DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS images.images;
DROP SCHEMA IF EXISTS images;