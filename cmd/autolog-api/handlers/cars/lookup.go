package cars

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	autologjwt "github.com/keola-dunn/autolog/cmd/autolog-api/jwt"
	"github.com/keola-dunn/autolog/internal/httputil"
	"github.com/keola-dunn/autolog/internal/logger"
	"github.com/keola-dunn/autolog/internal/service/car"
)

type lookupRequestParams struct {
	VIN string

	CarId string

	// PlateNumber string
	// State       string
}

type lookupResponsePlate struct {
	Number string `json:"number"`
	State  string `json:"state"`
}

type lookupResponse struct {
	///////////////
	// from cars db tables
	///////////////
	VIN          string               `json:"vin"`
	LicensePlate *lookupResponsePlate `json:"licensePlate,omitempty"`
	Year         int64                `json:"year"`
	Make         string               `json:"make"`
	Model        string               `json:"model"`
	Color        string               `json:"color"`

	//TransmissionStyle string `json:"transmissionStyle"`

	///////////////
	// from NHTSA VPIC
	///////////////
	ManufactureCity    string `json:"manufactureCity"`
	ManufactureState   string `json:"manufactureState"`
	ManufactureCountry string `json:"manufactureCountry"`

	///////////////
	// from service records tables
	///////////////

}

type lookupResponseAuthenticated struct {
	lookupResponse
	OwnerId string `json:"ownerId"`
}

// Lookup is the search request for cars. Can be searched for by Public autolog ID of the
// car, VIN, or plate.
// Will search autolog for car record - will also search NHTSA for details about the car
// Uses for this endpoint: finding a car to add to your garage, or looking up a car to
// see it's history. Presumably, using this, you aren't looking for your own car, you're looking
// at someone elses car.
// Public/Private responses are different.
func (h *CarsHandler) Lookup(w http.ResponseWriter, r *http.Request) {
	logEntry := logger.GetLogEntry(r)
	var userId string

	authToken := autologjwt.GetTokenFromAuthHeader(r.Header.Get("Authorization"))
	if strings.TrimSpace(authToken) != "" {
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

		userId = token.GetUserId()
	}

	var queryParams = make(url.Values, len(r.URL.Query()))
	for key, val := range r.URL.Query() {
		// convert all keys to lower case for ease of use
		queryParams[strings.ToLower(key)] = val
	}

	vin := strings.TrimSpace(queryParams.Get("vin"))
	carId := strings.TrimSpace(queryParams.Get("carid"))
	id := strings.TrimSpace(queryParams.Get("id"))
	// plateNumber := strings.TrimSpace(queryParams.Get("platenumber"))
	// state := strings.TrimSpace(queryParams.Get("state"))

	// decodeVINOutput, err := h.nhtsaClient.DecodeVINFlat(r.Context(), nhtsavpic.DecodeVINFlatInput{
	// 	VIN: vin,
	// })
	// if err != nil {
	// 	logEntry.Error("failed to decode vin", err)
	// 	httputil.RespondWithError(w, http.StatusInternalServerError, "")
	// 	return
	// }

	// if decodeVINOutput.Count <= 0 {
	// 	logEntry.Error("vin not found")
	// }

	getCarOutput, err := h.carService.GetCar(r.Context(), car.GetCarInput{
		VIN:      vin,
		PublicId: carId,
		Id:       id,
	})
	if err != nil {
		if errors.Is(err, car.ErrNotFound) {
			httputil.RespondWithError(w, http.StatusNotFound, "")
			return
		}
		logEntry.Error("failed to get car", err)
		httputil.RespondWithError(w, http.StatusInternalServerError, "")
		return
	}

	response := lookupResponse{
		VIN:          getCarOutput.VIN,
		LicensePlate: nil,
		Year:         getCarOutput.Year,
		Make:         getCarOutput.Make,
		Model:        getCarOutput.Model,
		Color:        getCarOutput.Color,
	}

	httputil.RespondWithJSON(w, http.StatusOK, response)
}
