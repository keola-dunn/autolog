package car

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type NHTSAVPICData struct {
	id                      string
	carId                   string
	VIN                     string
	Make                    string
	Model                   string
	Year                    int64
	Trim                    string
	Trim2                   string
	Manufacturer            string
	ManufacturerId          string
	PlantCompanyName        string
	PlantCity               string
	PlantState              string
	PlantCountry            string
	DisplacementCubicInches string
	DisplacementLiters      string
	DriveType               string
	EngineConfiguration     string
	EngineCylinders         string
	EngineHP                string
	EngineKW                string
	EngineManufacturer      string
	EngineModel             string
	FuelTypePrimary         string
	FuelTypeSecondary       string
	GCWR                    string
	GVWR                    string
	Seats                   string
	SeatsRows               string
	SteeringLocation        string
	TransmissionStyle       string
	TransmissionSpeeds      string
	VehicleType             string
	ValveTrainDesign        string
	WheelbaseLong           string
	WheelbaseShort          string
	WheelbaseType           string
	WheelSizeFront          string
	WheelSizeRear           string
	Payload                 []byte
	createdAt               time.Time
	updatedAt               time.Time
}

func (n *NHTSAVPICData) Id() string {
	return n.id
}

func (n *NHTSAVPICData) CarId() string {
	return n.carId
}

func (n *NHTSAVPICData) CreatedAt() time.Time {
	return n.createdAt
}

func (n *NHTSAVPICData) UpdatedAt() time.Time {
	return n.updatedAt
}

func createNHTSAVPICDataRecord(ctx context.Context, tx pgx.Tx, input NHTSAVPICData) error {
	query := `
	INSERT INTO nhtsa_vpic_data(
    	car_id,
    	vin,
    	make, 
    	model,
    	year,
    	trim,
    	trim2,
    	manufacturer,
		manufacturer_id,
    	plant_company_name,
    	plant_city,
    	plant_state,
    	plant_country,
    	displacement_ci,
    	displacement_l,
    	drive_type,
    	engine_configuration,
    	engine_cylinders,
    	engine_hp,
    	engine_kw,
    	engine_manufacturer,
    	engine_model,
    	fuel_type_primary,
    	fuel_type_secondary,
    	gcwr,
    	gvwr,
    	seats,
    	seats_rows,
    	steering_location,
    	transmission_style,
    	transmission_speeds,
    	vehicle_type,
    	valve_train_design,
		wheel_base_long,
		wheel_base_short,
		wheel_base_type, 
		wheel_size_front,
		wheel_size_rear,
    	payload
	) VALUES (
		$1,
		$2,
		$3,
		$4,
		$5,
		$6,
		$7,
		$8,
		$9,
		$10,
		$11,
		$12,
		$13,
		$14,
		$15,
		$16,
		$17,
		$18,
		$19,
		$20,
		$21,
		$22,
		$23,
		$24,
		$25,
		$26,
		$27,
		$28,
		$29,
		$30,
		$31,
		$32,
		$33,
		$34,
		$35,
		$36,
		$37,
		$38,
		$39)`

	if _, err := tx.Exec(ctx, query,
		input.CarId,
		input.VIN,
		input.Make,
		input.Model,
		input.Year,
		input.Trim,
		input.Trim2,
		input.Manufacturer,
		input.ManufacturerId,
		input.PlantCompanyName,
		input.PlantCity,
		input.PlantState,
		input.PlantCountry,
		input.DisplacementCubicInches,
		input.DisplacementLiters,
		input.DriveType,
		input.EngineConfiguration,
		input.EngineCylinders,
		input.EngineHP,
		input.EngineKW,
		input.EngineManufacturer,
		input.EngineModel,
		input.FuelTypePrimary,
		input.FuelTypeSecondary,
		input.GCWR,
		input.GVWR,
		input.Seats,
		input.SeatsRows,
		input.SteeringLocation,
		input.TransmissionStyle,
		input.TransmissionSpeeds,
		input.VehicleType,
		input.ValveTrainDesign,
		input.WheelbaseLong,
		input.WheelbaseShort,
		input.WheelbaseType,
		input.WheelSizeFront,
		input.WheelSizeRear,
		input.Payload); err != nil {
		return fmt.Errorf("failed to exec insert query: %w", err)
	}

	return nil
}
