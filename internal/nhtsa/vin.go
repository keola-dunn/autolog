package nhtsavpic

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// DecodeVINInput contains the input values needed/recommended to decode a VIN.
type DecodeVINInput struct {
	// VIN to decode. The VIN cannot be empty. The VIN does
	// not need to be 17 characters. Partial VINs and VINs pre-1980 can be submitted. Unknown characters
	// in a partial VIN can be replaced with "*"
	VIN string `json:"VIN"`

	// ModelYear is the year of the vehicle for the VIN to be decoded. This is not required.
	ModelYear int `json:"ModelYear"`
}

// DecodeVINOutput contains the output values from decoding a VIN.
type DecodeVINOutput struct {
	Count          int               `json:"Count"`
	Message        string            `json:"Message"`
	SearchCriteria string            `json:"SearchCriteria"`
	Results        []DecodeVINResult `json:"Results"`
}

// DecodeVINResult is the generic result for information retrieved about a VIN. It
// is a simple key-value pair data structure.
type DecodeVINResult struct {
	// Value is the value of the data type for the VIN decode result
	Value   string `json:"Value"`
	ValueID string `json:"ValueId"`

	// Variable is the key for the VIN decode result
	Variable   string `json:"Variable"`
	VariableID int    `json:"VariableId"`
}

// DecodeVIN retrieves any information NHTSA contains about the provided VIN.
func (c *Client) DecodeVIN(ctx context.Context, in DecodeVINInput) (DecodeVINOutput, error) {
	if strings.TrimSpace(in.VIN) == "" {
		return DecodeVINOutput{}, ErrInvalidArgument
	}

	var query string = "format=json"
	if in.ModelYear != 0 {
		query = fmt.Sprint(query, "&modelyear=", in.ModelYear)
	}

	url := url.URL{
		Scheme:   "https",
		Host:     baseURL,
		Path:     fmt.Sprintf("api/vehicles/decodevin/%s", strings.TrimSpace(in.VIN)),
		RawQuery: query,
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		return DecodeVINOutput{}, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return DecodeVINOutput{}, fmt.Errorf("failed to get models for make: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return DecodeVINOutput{}, fmt.Errorf("failed to decode VIN with status code: %d", resp.StatusCode)
	}

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return DecodeVINOutput{}, fmt.Errorf("failed to read response: %w", err)
	}

	var out DecodeVINOutput
	if err := json.Unmarshal(respData, &out); err != nil {
		return DecodeVINOutput{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return out, nil
}

type DecodeVINFlatInput struct {
	VIN       string `json:"VIN"`
	ModelYear int    `json:"ModelYear"`
}

type DecodeVINFlatOutput struct {
	Count          int                   `json:"Count"`
	Message        string                `json:"Message"`
	SearchCriteria string                `json:"SearchCriteria"`
	Results        []DecodeVINFlatResult `json:"Results"`
}

type DecodeVINFlatResult struct {
	ABS                                 string `json:"ABS"`
	ActiveSafetySysNote                 string `json:"ActiveSafetySysNote"`
	AdaptiveCruiseControl               string `json:"AdaptiveCruiseControl"`
	AdaptiveDrivingBeam                 string `json:"AdaptiveDrivingBeam"`
	AdaptiveHeadlights                  string `json:"AdaptiveHeadlights"`
	AdditionalErrorText                 string `json:"AdditionalErrorText"`
	AirBagLocCurtain                    string `json:"AirBagLocCurtain"`
	AirBagLocFront                      string `json:"AirBagLocFront"`
	AirBagLocKnee                       string `json:"AirBagLocKnee"`
	AirBagLocSeatCushion                string `json:"AirBagLocSeatCushion"`
	AirBagLocSide                       string `json:"AirBagLocSide"`
	AutoReverseSystem                   string `json:"AutoReverseSystem"`
	AutomaticPedestrianAlertingSound    string `json:"AutomaticPedestrianAlertingSound"`
	AxleConfiguration                   string `json:"AxleConfiguration"`
	Axles                               string `json:"Axles"`
	BasePrice                           string `json:"BasePrice"`
	BatteryA                            string `json:"BatteryA"`
	BatteryA_to                         string `json:"BatteryA_to"`
	BatteryCells                        string `json:"BatteryCells"`
	BatteryInfo                         string `json:"BatteryInfo"`
	BatteryKWh                          string `json:"BatteryKWh"`
	BatteryKWh_to                       string `json:"BatteryKWh_to"`
	BatteryModules                      string `json:"BatteryModules"`
	BatteryPacks                        string `json:"BatteryPacks"`
	BatteryType                         string `json:"BatteryType"`
	BatteryV                            string `json:"BatteryV"`
	BatteryV_to                         string `json:"BatteryV_to"`
	BedLengthIN                         string `json:"BedLengthIN"`
	BedType                             string `json:"BedType"`
	BlindSpotMon                        string `json:"BlindSpotMon"`
	BodyCabType                         string `json:"BodyCabType"`
	BodyClass                           string `json:"BodyClass"`
	BrakeSystemDesc                     string `json:"BrakeSystemDesc"`
	BrakeSystemType                     string `json:"BrakeSystemType"`
	BusFloorConfigType                  string `json:"BusFloorConfigType"`
	BusLength                           string `json:"BusLength"`
	BusType                             string `json:"BusType"`
	CAN_AACN                            string `json:"CAN_AACN"`
	CIB                                 string `json:"CIB"`
	CashForClunkers                     string `json:"CashForClunkers"`
	ChargerLevel                        string `json:"ChargerLevel"`
	ChargerPowerKW                      string `json:"ChargerPowerKW"`
	CoolingType                         string `json:"CoolingType"`
	CurbWeightLB                        string `json:"CurbWeightLB"`
	CustomMotorcycleType                string `json:"CustomMotorcycleType"`
	DaytimeRunningLight                 string `json:"DaytimeRunningLight"`
	DestinationMarket                   string `json:"DestinationMarket"`
	DisplacementCC                      string `json:"DisplacementCC"`
	DisplacementCI                      string `json:"DisplacementCI"`
	DisplacementL                       string `json:"DisplacementL"`
	Doors                               string `json:"Doors"`
	DriveType                           string `json:"DriveType"`
	DriverAssist                        string `json:"DriverAssist"`
	DynamicBrakeSupport                 string `json:"DynamicBrakeSupport"`
	EDR                                 string `json:"EDR"`
	ESC                                 string `json:"ESC"`
	EVDriveUnit                         string `json:"EVDriveUnit"`
	ElectrificationLevel                string `json:"ElectrificationLevel"`
	EngineConfiguration                 string `json:"EngineConfiguration"`
	EngineCycles                        string `json:"EngineCycles"`
	EngineCylinders                     string `json:"EngineCylinders"`
	EngineHP                            string `json:"EngineHP"`
	EngineHP_to                         string `json:"EngineHP_to"`
	EngineKW                            string `json:"EngineKW"`
	EngineManufacturer                  string `json:"EngineManufacturer"`
	EngineModel                         string `json:"EngineModel"`
	EntertainmentSystem                 string `json:"EntertainmentSystem"`
	ErrorCode                           string `json:"ErrorCode"`
	ErrorText                           string `json:"ErrorText"`
	ForwardCollisionWarning             string `json:"ForwardCollisionWarning"`
	FuelInjectionType                   string `json:"FuelInjectionType"`
	FuelTypePrimary                     string `json:"FuelTypePrimary"`
	FuelTypeSecondary                   string `json:"FuelTypeSecondary"`
	GCWR                                string `json:"GCWR"`
	GCWR_to                             string `json:"GCWR_to"`
	GVWR                                string `json:"GVWR"`
	GVWR_to                             string `json:"GVWR_to"`
	KeylessIgnition                     string `json:"KeylessIgnition"`
	LaneDepartureWarning                string `json:"LaneDepartureWarning"`
	LaneKeepSystem                      string `json:"LaneKeepSystem"`
	LowerBeamHeadlampLightSource        string `json:"LowerBeamHeadlampLightSource"`
	Make                                string `json:"Make"`
	MakeID                              string `json:"MakeID"`
	Manufacturer                        string `json:"Manufacturer"`
	ManufacturerId                      string `json:"ManufacturerId"`
	Model                               string `json:"Model"`
	ModelID                             string `json:"ModelID"`
	ModelYear                           string `json:"ModelYear"`
	MotorcycleChassisType               string `json:"MotorcycleChassisType"`
	MotorcycleSuspensionType            string `json:"MotorcycleSuspensionType"`
	NCSABodyType                        string `json:"NCSABodyType"`
	NCSAMake                            string `json:"NCSAMake"`
	NCSAMapExcApprovedBy                string `json:"NCSAMapExcApprovedBy"`
	NCSAMapExcApprovedOn                string `json:"NCSAMapExcApprovedOn"`
	NCSAMappingException                string `json:"NCSAMappingException"`
	NCSAModel                           string `json:"NCSAModel"`
	NCSANote                            string `json:"NCSANote"`
	Note                                string `json:"Note"`
	OtherBusInfo                        string `json:"OtherBusInfo"`
	OtherEngineInfo                     string `json:"OtherEngineInfo"`
	OtherMotorcycleInfo                 string `json:"OtherMotorcycleInfo"`
	OtherRestraintSystemInfo            string `json:"OtherRestraintSystemInfo"`
	OtherTrailerInfo                    string `json:"OtherTrailerInfo"`
	ParkAssist                          string `json:"ParkAssist"`
	PedestrianAutomaticEmergencyBraking string `json:"PedestrianAutomaticEmergencyBraking"`
	PlantCity                           string `json:"PlantCity"`
	PlantCompanyName                    string `json:"PlantCompanyName"`
	PlantCountry                        string `json:"PlantCountry"`
	PlantState                          string `json:"PlantState"`
	PossibleValues                      string `json:"PossibleValues"`
	Pretensioner                        string `json:"Pretensioner"`
	RearCrossTrafficAlert               string `json:"RearCrossTrafficAlert"`
	RearVisibilitySystem                string `json:"RearVisibilitySystem"`
	SAEAutomationLevel                  string `json:"SAEAutomationLevel"`
	SAEAutomationLevel_to               string `json:"SAEAutomationLevel_to"`
	SeatBeltsAll                        string `json:"SeatBeltsAll"`
	SeatRows                            string `json:"SeatRows"`
	Seats                               string `json:"Seats"`
	SemiautomaticHeadlampBeamSwitching  string `json:"SemiautomaticHeadlampBeamSwitching"`
	Series                              string `json:"Series"`
	Series2                             string `json:"Series2"`
	SteeringLocation                    string `json:"SteeringLocation"`
	SuggestedVIN                        string `json:"SuggestedVIN"`
	TPMS                                string `json:"TPMS"`
	TopSpeedMPH                         string `json:"TopSpeedMPH"`
	TrackWidth                          string `json:"TrackWidth"`
	TractionControl                     string `json:"TractionControl"`
	TrailerBodyType                     string `json:"TrailerBodyType"`
	TrailerLength                       string `json:"TrailerLength"`
	TrailerType                         string `json:"TrailerType"`
	TransmissionSpeeds                  string `json:"TransmissionSpeeds"`
	TransmissionStyle                   string `json:"TransmissionStyle"`
	Trim                                string `json:"Trim"`
	Trim2                               string `json:"Trim2"`
	Turbo                               string `json:"Turbo"`
	VIN                                 string `json:"VIN"`
	ValveTrainDesign                    string `json:"ValveTrainDesign"`
	VehicleType                         string `json:"VehicleType"`
	WheelBaseLong                       string `json:"WheelBaseLong"`
	WheelBaseShort                      string `json:"WheelBaseShort"`
	WheelBaseType                       string `json:"WheelBaseType"`
	WheelSizeFront                      string `json:"WheelSizeFront"`
	WheelSizeRear                       string `json:"WheelSizeRear"`
	Wheels                              string `json:"Wheels"`
	Windows                             string `json:"Windows"`
}

func (c *Client) DecodeVINFlat(ctx context.Context, in DecodeVINFlatInput) (DecodeVINFlatOutput, error) {

	// TODO: Should I filter out vins longer than 17 characters? I think that's standard length in the US
	if strings.TrimSpace(in.VIN) == "" {
		return DecodeVINFlatOutput{}, ErrInvalidArgument
	}

	var query string = "format=json"
	if in.ModelYear > 1900 {
		query = fmt.Sprint(query, "&modelyear=", in.ModelYear)
	}

	url := url.URL{
		Scheme:   "https",
		Host:     baseURL,
		Path:     fmt.Sprintf("api/vehicles/decodevinvalues/%s", strings.TrimSpace(in.VIN)),
		RawQuery: query,
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		return DecodeVINFlatOutput{}, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return DecodeVINFlatOutput{}, fmt.Errorf("failed to get models for make: %w", err)
	}

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return DecodeVINFlatOutput{}, fmt.Errorf("failed to read response: %w", err)
	}

	var out DecodeVINFlatOutput
	if err := json.Unmarshal(respData, &out); err != nil {
		return DecodeVINFlatOutput{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return out, nil
}

type DecodeVINExtendedInput struct {
	VIN       string `json:"VIN"`
	ModelYear int    `json:"ModelYear"`
}

type DecodeVINExtendedOutput struct {
	Count          int               `json:"Count"`
	Message        string            `json:"Message"`
	SearchCriteria string            `json:"SearchCriteria"`
	Results        []DecodeVINResult `json:"Results"`
}

func (c *Client) DecodeVINExtended(ctx context.Context, in DecodeVINExtendedInput) (DecodeVINExtendedOutput, error) {
	// TODO: Should I filter out vins longer than 17 characters? I think that's standard length in the US
	if strings.TrimSpace(in.VIN) == "" {
		return DecodeVINExtendedOutput{}, ErrInvalidArgument
	}

	var query string = "format=json"
	if in.ModelYear > 1900 {
		query = fmt.Sprint(query, "&modelyear=", in.ModelYear)
	}

	url := url.URL{
		Scheme:   "https",
		Host:     baseURL,
		Path:     fmt.Sprintf("api/vehicles/decodevinextended/%s", strings.TrimSpace(in.VIN)),
		RawQuery: query,
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		return DecodeVINExtendedOutput{}, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return DecodeVINExtendedOutput{}, fmt.Errorf("failed to do request: %w", err)
	}

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return DecodeVINExtendedOutput{}, fmt.Errorf("failed to read response: %w", err)
	}

	var out DecodeVINExtendedOutput
	if err := json.Unmarshal(respData, &out); err != nil {
		return DecodeVINExtendedOutput{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return out, nil
}

type DecodeVINExtendedFlatInput struct {
	VIN       string `json:"VIN"`
	ModelYear int    `json:"ModelYear"`
}

type DecodeVINExtendedFlatOutput struct {
	Count          int                            `json:"Count"`
	Message        string                         `json:"Message"`
	SearchCriteria string                         `json:"SearchCriteria"`
	Results        []DecodeVINExtendedFlatResults `json:"Results"`
}

type DecodeVINExtendedFlatResults struct {
	ABS                                 string `json:"ABS"`
	ActiveSafetySysNote                 string `json:"ActiveSafetySysNote"`
	AdaptiveCruiseControl               string `json:"AdaptiveCruiseControl"`
	AdaptiveDrivingBeam                 string `json:"AdaptiveDrivingBeam"`
	AdaptiveHeadlights                  string `json:"AdaptiveHeadlights"`
	AdditionalErrorText                 string `json:"AdditionalErrorText"`
	AirBagLocCurtain                    string `json:"AirBagLocCurtain"`
	AirBagLocFront                      string `json:"AirBagLocFront"`
	AirBagLocKnee                       string `json:"AirBagLocKnee"`
	AirBagLocSeatCushion                string `json:"AirBagLocSeatCushion"`
	AirBagLocSide                       string `json:"AirBagLocSide"`
	AutoReverseSystem                   string `json:"AutoReverseSystem"`
	AutomaticPedestrianAlertingSound    string `json:"AutomaticPedestrianAlertingSound"`
	AxleConfiguration                   string `json:"AxleConfiguration"`
	Axles                               string `json:"Axles"`
	BasePrice                           string `json:"BasePrice"`
	BatteryA                            string `json:"BatteryA"`
	BatteryA_to                         string `json:"BatteryA_to"`
	BatteryCells                        string `json:"BatteryCells"`
	BatteryInfo                         string `json:"BatteryInfo"`
	BatteryKWh                          string `json:"BatteryKWh"`
	BatteryKWh_to                       string `json:"BatteryKWh_to"`
	BatteryModules                      string `json:"BatteryModules"`
	BatteryPacks                        string `json:"BatteryPacks"`
	BatteryType                         string `json:"BatteryType"`
	BatteryV                            string `json:"BatteryV"`
	BatteryV_to                         string `json:"BatteryV_to"`
	BedLengthIN                         string `json:"BedLengthIN"`
	BedType                             string `json:"BedType"`
	BlindSpotMon                        string `json:"BlindSpotMon"`
	BodyCabType                         string `json:"BodyCabType"`
	BodyClass                           string `json:"BodyClass"`
	BrakeSystemDesc                     string `json:"BrakeSystemDesc"`
	BrakeSystemType                     string `json:"BrakeSystemType"`
	BusFloorConfigType                  string `json:"BusFloorConfigType"`
	BusLength                           string `json:"BusLength"`
	BusType                             string `json:"BusType"`
	CAN_AACN                            string `json:"CAN_AACN"`
	CIB                                 string `json:"CIB"`
	CashForClunkers                     string `json:"CashForClunkers"`
	ChargerLevel                        string `json:"ChargerLevel"`
	ChargerPowerKW                      string `json:"ChargerPowerKW"`
	CoolingType                         string `json:"CoolingType"`
	CurbWeightLB                        string `json:"CurbWeightLB"`
	CustomMotorcycleType                string `json:"CustomMotorcycleType"`
	DaytimeRunningLight                 string `json:"DaytimeRunningLight"`
	DestinationMarket                   string `json:"DestinationMarket"`
	DisplacementCC                      string `json:"DisplacementCC"`
	DisplacementCI                      string `json:"DisplacementCI"`
	DisplacementL                       string `json:"DisplacementL"`
	Doors                               string `json:"Doors"`
	DriveType                           string `json:"DriveType"`
	DriverAssist                        string `json:"DriverAssist"`
	DynamicBrakeSupport                 string `json:"DynamicBrakeSupport"`
	EDR                                 string `json:"EDR"`
	ESC                                 string `json:"ESC"`
	EVDriveUnit                         string `json:"EVDriveUnit"`
	ElectrificationLevel                string `json:"ElectrificationLevel"`
	EngineConfiguration                 string `json:"EngineConfiguration"`
	EngineCycles                        string `json:"EngineCycles"`
	EngineCylinders                     string `json:"EngineCylinders"`
	EngineHP                            string `json:"EngineHP"`
	EngineHP_to                         string `json:"EngineHP_to"`
	EngineKW                            string `json:"EngineKW"`
	EngineManufacturer                  string `json:"EngineManufacturer"`
	EngineModel                         string `json:"EngineModel"`
	EntertainmentSystem                 string `json:"EntertainmentSystem"`
	ErrorCode                           string `json:"ErrorCode"`
	ErrorText                           string `json:"ErrorText"`
	ForwardCollisionWarning             string `json:"ForwardCollisionWarning"`
	FuelInjectionType                   string `json:"FuelInjectionType"`
	FuelTypePrimary                     string `json:"FuelTypePrimary"`
	FuelTypeSecondary                   string `json:"FuelTypeSecondary"`
	GCWR                                string `json:"GCWR"`
	GCWR_to                             string `json:"GCWR_to"`
	GVWR                                string `json:"GVWR"`
	GVWR_to                             string `json:"GVWR_to"`
	KeylessIgnition                     string `json:"KeylessIgnition"`
	LaneDepartureWarning                string `json:"LaneDepartureWarning"`
	LaneKeepSystem                      string `json:"LaneKeepSystem"`
	LowerBeamHeadlampLightSource        string `json:"LowerBeamHeadlampLightSource"`
	Make                                string `json:"Make"`
	MakeID                              string `json:"MakeID"`
	Manufacturer                        string `json:"Manufacturer"`
	ManufacturerId                      string `json:"ManufacturerId"`
	Model                               string `json:"Model"`
	ModelID                             string `json:"ModelID"`
	ModelYear                           string `json:"ModelYear"`
	MotorcycleChassisType               string `json:"MotorcycleChassisType"`
	MotorcycleSuspensionType            string `json:"MotorcycleSuspensionType"`
	NCSABodyType                        string `json:"NCSABodyType"`
	NCSAMake                            string `json:"NCSAMake"`
	NCSAMapExcApprovedBy                string `json:"NCSAMapExcApprovedBy"`
	NCSAMapExcApprovedOn                string `json:"NCSAMapExcApprovedOn"`
	NCSAMappingException                string `json:"NCSAMappingException"`
	NCSAModel                           string `json:"NCSAModel"`
	NCSANote                            string `json:"NCSANote"`
	Note                                string `json:"Note"`
	OtherBusInfo                        string `json:"OtherBusInfo"`
	OtherEngineInfo                     string `json:"OtherEngineInfo"`
	OtherMotorcycleInfo                 string `json:"OtherMotorcycleInfo"`
	OtherRestraintSystemInfo            string `json:"OtherRestraintSystemInfo"`
	OtherTrailerInfo                    string `json:"OtherTrailerInfo"`
	ParkAssist                          string `json:"ParkAssist"`
	PedestrianAutomaticEmergencyBraking string `json:"PedestrianAutomaticEmergencyBraking"`
	PlantCity                           string `json:"PlantCity"`
	PlantCompanyName                    string `json:"PlantCompanyName"`
	PlantCountry                        string `json:"PlantCountry"`
	PlantState                          string `json:"PlantState"`
	PossibleValues                      string `json:"PossibleValues"`
	Pretensioner                        string `json:"Pretensioner"`
	RearCrossTrafficAlert               string `json:"RearCrossTrafficAlert"`
	RearVisibilitySystem                string `json:"RearVisibilitySystem"`
	SAEAutomationLevel                  string `json:"SAEAutomationLevel"`
	SAEAutomationLevel_to               string `json:"SAEAutomationLevel_to"`
	SeatBeltsAll                        string `json:"SeatBeltsAll"`
	SeatRows                            string `json:"SeatRows"`
	Seats                               string `json:"Seats"`
	SemiautomaticHeadlampBeamSwitching  string `json:"SemiautomaticHeadlampBeamSwitching"`
	Series                              string `json:"Series"`
	Series2                             string `json:"Series2"`
	SteeringLocation                    string `json:"SteeringLocation"`
	SuggestedVIN                        string `json:"SuggestedVIN"`
	TPMS                                string `json:"TPMS"`
	TopSpeedMPH                         string `json:"TopSpeedMPH"`
	TrackWidth                          string `json:"TrackWidth"`
	TractionControl                     string `json:"TractionControl"`
	TrailerBodyType                     string `json:"TrailerBodyType"`
	TrailerLength                       string `json:"TrailerLength"`
	TrailerType                         string `json:"TrailerType"`
	TransmissionSpeeds                  string `json:"TransmissionSpeeds"`
	TransmissionStyle                   string `json:"TransmissionStyle"`
	Trim                                string `json:"Trim"`
	Trim2                               string `json:"Trim2"`
	Turbo                               string `json:"Turbo"`
	VIN                                 string `json:"VIN"`
	ValveTrainDesign                    string `json:"ValveTrainDesign"`
	VehicleType                         string `json:"VehicleType"`
	WheelBaseLong                       string `json:"WheelBaseLong"`
	WheelBaseShort                      string `json:"WheelBaseShort"`
	WheelBaseType                       string `json:"WheelBaseType"`
	WheelSizeFront                      string `json:"WheelSizeFront"`
	WheelSizeRear                       string `json:"WheelSizeRear"`
	Wheels                              string `json:"Wheels"`
	Windows                             string `json:"Windows"`
}

func (c *Client) DecodeVINExtendedFlat(ctx context.Context, in DecodeVINExtendedFlatInput) (DecodeVINExtendedFlatOutput, error) {
	if strings.TrimSpace(in.VIN) == "" {
		return DecodeVINExtendedFlatOutput{}, ErrInvalidArgument
	}

	var query string = "format=json"
	if in.ModelYear > 1900 {
		query = fmt.Sprint(query, "&modelyear=", in.ModelYear)
	}

	url := url.URL{
		Scheme:   "https",
		Host:     baseURL,
		Path:     fmt.Sprintf("api/vehicles/decodevinvaluesextended/%s", strings.TrimSpace(in.VIN)),
		RawQuery: query,
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		return DecodeVINExtendedFlatOutput{}, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return DecodeVINExtendedFlatOutput{}, fmt.Errorf("failed to do request: %w", err)
	}

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return DecodeVINExtendedFlatOutput{}, fmt.Errorf("failed to read response: %w", err)
	}

	var out DecodeVINExtendedFlatOutput
	if err := json.Unmarshal(respData, &out); err != nil {
		return DecodeVINExtendedFlatOutput{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return out, nil
}
