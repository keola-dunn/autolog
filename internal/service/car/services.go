package car

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type VehicleService interface {
	Name()

	Scan(value interface{}) error
	Value() (driver.Value, error)
}

func QuartsToLiters(quarts float64) float64 {
	return quarts / 1.057
}

func GallonsToLiters(gallons float64) float64 {
	return gallons * 3.785

}

type OilChangeService struct {
	OilBrand       string  `json:"brand"`
	Viscosity      string  `json:"viscosity"`
	VolumeLiters   float64 `json:"volumeLiters"`
	Filter         string  `json:"filter"`
	NewCrushWasher bool    `json:"newCrushWasher"`
}

func (*OilChangeService) Name() string {
	return "oil-change"
}

func (o OilChangeService) Value() (driver.Value, error) {
	return json.Marshal(o)
}

func (o *OilChangeService) Scan(value any) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &o)
}

type TirePosition string

const (
	TirePositionLeftFront  = TirePosition("LF")
	TirePositionRightFront = TirePosition("RF")
	TirePositionLeftRear   = TirePosition("LR")
	TirePositionRightRear  = TirePosition("RR")

	TirePostionUnknown = TirePosition("unknown")
)

// UnmarshalJSON implements the json.Unmarshaler interface for TirePosition.
func (t *TirePosition) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("failed to unmarshal JSON into string: %w", err)
	}

	var tp TirePosition
	switch s {
	case "LF":
		tp = TirePositionLeftFront
	case "RF":
		tp = TirePositionRightFront
	case "LR":
		tp = TirePositionLeftRear
	case "RR":
		tp = TirePositionRightRear
	default:
	}

	t = &tp
	return nil
}

func (t TirePosition) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(t))
}

type TireChangeService struct {
	TireBrand string `json:"brand"`
	Tire      string `json:"tire"`
	FrontSize string `json:"frontSize"`
	RearSize  string `json:"rearSize"`

	TiresChanged []TirePosition `json:"tiresChanged,omitempty"`
	IsDually     bool           `json:"isDually"`
}

func (*TireChangeService) Name() string {
	return "tire-change"
}

func (t TireChangeService) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *TireChangeService) Scan(value any) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &t)
}

// CoolantType denotes the typical Coolant types a car will use
// Data source is here: https://www.valvolineglobal.com/en/engine-coolant/
type CoolantType string

const (
	// CoolantTypeIAT - Inorganic Additive Technology. Uses Silicates as the inhibitor. Typically
	// green. Common in old cars.
	CoolantTypeIAT = CoolantType("IAT")

	// CoolantTypeOAT - Organic Acid Technology. Uses organic acids as the inhibitors. Typically
	// orange. Also referred to as Dexcool. Common in GM + VW.
	CoolantTypeOAT = CoolantType("OAT")

	// CoolantTypeHOAT - Hybrid Organic Acid Technology. Silicates and organic acids as inhibitors.
	// Typically yellow. Ford, Chrysler, some euros.
	CoolantTypeHOAT = CoolantType("HOAT")

	// CoolantTypePhosphateFreeHOAT - Phosphate Free Hybrid Organic Acid Technology. Commonly
	// turquoise. Tesla, BMW, Mini, Volvo
	CoolantTypePhosphateFreeHOAT = CoolantType("Phosphate Free HOAT")

	// CoolantTypePHOAT - Phosphated Hybrid Organic Acid Technology. Blue or Pink. Toyota, Nissan,
	// Subaru, asian vehicles.
	CoolantTypePHOAT = CoolantType("PHOAT")

	// CoolantTypeSiHOAT - Silicated Hybrid Organic Acid Technology. Purple. Porsche, Mercedes-Benz,
	// Audi, etc.
	CoolantTypeSiHOAT = CoolantType("Si-HOAT")
)

type CoolantFlushService struct {
	CoolantBrand string      `json:"brand"`
	CoolantColor string      `json:"color"`
	CoolantType  CoolantType `json:"type"`
	VolumeLiters float64     `json:"volumeLiters"`
}

func (*CoolantFlushService) Name() string {
	return "coolant-flush"
}

func (c CoolantFlushService) Value() (driver.Value, error) {
	return json.Marshal(c)
}

func (c *CoolantFlushService) Scan(value any) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &c)
}
