-- +goose Up

CREATE TABLE IF NOT EXISTS nhtsa_vpic_data (
    id uuid NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    car_id uuid NOT NULL references cars(id),
    vin varchar(32),
    make varchar(64), 
    model varchar(64),
    year varchar(8),
    trim varchar(64),
    trim2 varchar(64),
    manufacturer varchar(64),
	manufacturer_id varchar(64),
    plant_company_name varchar(64),
    plant_city varchar(256),
    plant_state varchar(256),
    plant_country varchar(256),
    displacement_ci varchar(64),
    displacement_l varchar(64),
    drive_type varchar(64),
    engine_configuration varchar(64),
    engine_cylinders varchar(32),
    engine_hp varchar(32),
    engine_kw varchar(32),
    engine_manufacturer varchar(128),
    engine_model varchar(64),
    fuel_type_primary varchar(64),
    fuel_type_secondary varchar(64),
    gcwr varchar(64),
    gvwr varchar(64),
    seats varchar(16),
    seats_rows varchar(8),
    steering_location varchar(32),
    transmission_style varchar(128),
    transmission_speeds varchar(16),
    vehicle_type varchar(32),
    valve_train_design varchar(32),
	wheel_base_long varchar(16),
	wheel_base_short varchar(16),
	wheel_base_type varchar(16), 
	wheel_size_front varchar(16),
	wheel_size_rear varchar(16),
    payload JSONB default '{}'::JSONB,
    created_at timestamptz DEFAULT NOW(),
    updated_at timestamptz DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_nhtsa_vpic_data_car_id ON nhtsa_vpic_data(car_id);
CREATE INDEX IF NOT EXISTS idx_nhtsa_vpic_data_make_model_year ON nhtsa_vpic_data(make, model, year);
CREATE INDEX IF NOT EXISTS idx_nhtsa_vpic_data_vin ON nhtsa_vpic_data(vin);

-- +goose Down
DROP TABLE IF EXISTS nhtsa_vpic_data;
