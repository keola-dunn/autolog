-- +goose Up
CREATE INDEX IF NOT EXISTS idx_service_logs_car_id ON service_logs(car_id);

-- +goose Down
DROP INDEX IF EXISTS idx_service_logs_car_id;
