-- +goose Up

ALTER TABLE cars ADD COLUMN IF NOT EXISTS public_id varchar(6) UNIQUE; 
ALTER TABLE cars ADD COLUMN IF NOT EXISTS transmission_style varchar(128); 
ALTER TABLE cars ADD COLUMN IF NOT EXISTS manufacture_city varchar(256); 
ALTER TABLE cars ADD COLUMN IF NOT EXISTS manufacture_state varchar(256); 
ALTER TABLE cars ADD COLUMN IF NOT EXISTS manufacture_country varchar(256); 

DisplacementCI                      string `json:"DisplacementCI"`
DisplacementL                       string `json:"DisplacementL"`

CREATE INDEX IF NOT EXISTS idx_cars_public_id ON cars(public_id);

-- +goose Down
DROP INDEX IF EXISTS idx_cars_public_id;

ALTER TABLE cars
DROP COLUMN IF EXISTS transmission_style;

ALTER TABLE cars
DROP COLUMN IF EXISTS manufacture_city;

ALTER TABLE cars
DROP COLUMN IF EXISTS manufacture_state;

ALTER TABLE cars
DROP COLUMN IF EXISTS manufacture_country;
