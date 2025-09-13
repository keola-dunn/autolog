package cars

import (
	"encoding/json"
	"io"
	"net/http"
	"slices"
	"strconv"

	"github.com/keola-dunn/autolog/internal/httputil"
	autologjwt "github.com/keola-dunn/autolog/internal/jwt"
	"github.com/keola-dunn/autolog/internal/logger"
	nhtsavpic "github.com/keola-dunn/autolog/internal/nhtsa"
	"github.com/keola-dunn/autolog/internal/service/car"
)

type createCarRequest struct {
	VIN   string `json:"vin"`
	Make  string `json:"make"`
	Model string `json:"model"`
	Year  int64  `json:"year"`
	Trim  string `json:"trim"`
	Color string `json:"color"`
}

//type createCarResponse struct {}

func (h *CarsHandler) CreateCar(w http.ResponseWriter, r *http.Request) {
	logEntry := logger.GetLogEntry(r)

	authToken := autologjwt.GetTokenFromAuthHeader(r.Header.Get("Authorization"))
	valid, token, err := autologjwt.VerifyToken(authToken, h.jwtSecret)
	if err != nil {
		logEntry.Error("failed to verify token", err)
		httputil.RespondWithError(w, http.StatusInternalServerError, "")
		return
	}
	if !valid {
		logEntry.Error("auth token is invalid", err)
		httputil.RespondWithError(w, http.StatusForbidden, "")
		return
	}

	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		logEntry.Error("failed to read request body", err)
		httputil.RespondWithError(w, http.StatusInternalServerError, "")
		return
	}

	var req createCarRequest
	if err := json.Unmarshal(requestBody, &req); err != nil {
		logEntry.Error("failed to unmarshal request body", err)
		httputil.RespondWithError(w, http.StatusBadRequest, "")
		return
	}

	decodedVINData, err := h.nhtsaClient.DecodeVINFlat(r.Context(), nhtsavpic.DecodeVINFlatInput{
		VIN:       req.VIN,
		ModelYear: int(req.Year),
	})
	if err != nil {
		logEntry.Error("failed to decode vin", err)
		httputil.RespondWithError(w, http.StatusInternalServerError, "")
		return
	}

	if decodedVINData.Count <= 0 || decodedVINData.Count > 1 {
		logEntry.Warn("car not found", "vin", req.VIN, "modelYear", req.Year)
		httputil.RespondWithError(w, http.StatusNotFound, "vin not found")
		return
	}

	errorCodes, err := decodedVINData.Results[0].ErrorCodes()
	if err != nil {
		logEntry.Error("failed to get error codes for decoded vin", err)
		httputil.RespondWithError(w, http.StatusInternalServerError, "")
		return
	}

	if !slices.Contains(errorCodes, nhtsavpic.ErrorCodeSuccess) {
		logEntry.Warn("nhtsavpic response doesn't indicate successful decode",
			"vin", req.VIN, "modelYear", req.Year)
		httputil.RespondWithError(w, http.StatusNotFound, "vin not found")
		return
	}

	modelYear, _ := strconv.Atoi(decodedVINData.Results[0].ModelYear)
	payload, _ := json.Marshal(decodedVINData.Results[0])

	if err := h.carService.CreateCar(r.Context(), token.GetUserId(), car.Car{
		Make:  req.Make,
		Model: req.Model,
		Year:  req.Year,
		Trim:  req.Trim,
		VIN:   req.VIN,
		Color: req.Color,
	}, car.NHTSAVPICData{
		VIN:                     decodedVINData.Results[0].VIN,
		Make:                    decodedVINData.Results[0].Make,
		Model:                   decodedVINData.Results[0].Model,
		Year:                    int64(modelYear),
		Trim:                    decodedVINData.Results[0].Trim,
		Trim2:                   decodedVINData.Results[0].Trim2,
		Manufacturer:            decodedVINData.Results[0].Manufacturer,
		ManufacturerId:          decodedVINData.Results[0].ManufacturerId,
		PlantCompanyName:        decodedVINData.Results[0].PlantCompanyName,
		PlantCity:               decodedVINData.Results[0].PlantCity,
		PlantState:              decodedVINData.Results[0].PlantState,
		PlantCountry:            decodedVINData.Results[0].PlantCountry,
		DisplacementCubicInches: decodedVINData.Results[0].DisplacementCI,
		DisplacementLiters:      decodedVINData.Results[0].DisplacementL,
		DriveType:               decodedVINData.Results[0].DriveType,
		EngineConfiguration:     decodedVINData.Results[0].EngineConfiguration,
		EngineCylinders:         decodedVINData.Results[0].EngineCylinders,
		EngineHP:                decodedVINData.Results[0].EngineHP,
		EngineKW:                decodedVINData.Results[0].EngineKW,
		EngineManufacturer:      decodedVINData.Results[0].EngineManufacturer,
		EngineModel:             decodedVINData.Results[0].EngineModel,
		FuelTypePrimary:         decodedVINData.Results[0].FuelTypePrimary,
		FuelTypeSecondary:       decodedVINData.Results[0].FuelTypeSecondary,
		GCWR:                    decodedVINData.Results[0].GCWR,
		GVWR:                    decodedVINData.Results[0].GVWR,
		Seats:                   decodedVINData.Results[0].Seats,
		SeatsRows:               decodedVINData.Results[0].SeatRows,
		SteeringLocation:        decodedVINData.Results[0].SteeringLocation,
		TransmissionStyle:       decodedVINData.Results[0].TransmissionStyle,
		TransmissionSpeeds:      decodedVINData.Results[0].TransmissionSpeeds,
		VehicleType:             decodedVINData.Results[0].VehicleType,
		ValveTrainDesign:        decodedVINData.Results[0].ValveTrainDesign,
		WheelbaseLong:           decodedVINData.Results[0].WheelBaseLong,
		WheelbaseShort:          decodedVINData.Results[0].WheelBaseShort,
		WheelbaseType:           decodedVINData.Results[0].WheelBaseType,
		WheelSizeFront:          decodedVINData.Results[0].WheelSizeFront,
		WheelSizeRear:           decodedVINData.Results[0].WheelSizeRear,
		Payload:                 payload}); err != nil {
		logEntry.Error("failed to create car", err)
		httputil.RespondWithError(w, http.StatusInternalServerError, "")
		return
	}

	httputil.RespondWithJSON(w, http.StatusCreated, req)
}
