-- +goose Up

ALTER TABLE cars ADD COLUMN IF NOT EXISTS color varchar(256);

CREATE TABLE IF NOT EXISTS license_plates (
    id uuid NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    plate_number varchar(8), 
    state varchar(2),
    country varchar(3) DEFAULT 'us',
    user_id uuid NOT NULL references users(id),
    created_at timestamptz DEFAULT NOW(),
    updated_at timestamptz DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_license_plates_plate_number ON license_plates(plate_number);
CREATE INDEX IF NOT EXISTS idx_license_plates_plate_state ON license_plates(plate_number, state);

-- +goose Down
ALTER TABLE cars DROP COLUMN IF EXISTS color;
DROP TABLE IF EXISTS license_plates;
