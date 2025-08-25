package cars

import (
	"encoding/json"
	"io"
	"net/http"

	autologjwt "github.com/keola-dunn/autolog/cmd/autolog-api/jwt"
	"github.com/keola-dunn/autolog/internal/httputil"
	"github.com/keola-dunn/autolog/internal/logger"
	"github.com/keola-dunn/autolog/internal/service/car"
)

type createCarRequest struct {
	VIN   string `json:"vin"`
	Make  string `json:"make"`
	Model string `json:"model"`
	Year  int64  `json:"year"`
	Trim  string `json:"trim"`
}

//type createCarResponse struct {}

func (h *CarsHandler) CreateCar(w http.ResponseWriter, r *http.Request) {
	logEntry := logger.GetLogEntry(r)

	authToken := autologjwt.GetTokenFromAuthHeader(r.Header.Get("Authorization"))
	ok, token, err := autologjwt.VerifyToken(authToken, h.jwtSecret)
	if err != nil {
		logEntry.Error("failed to verify token", err)
		httputil.RespondWithError(w, http.StatusInternalServerError, "")
		return
	}
	if !ok {
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
		httputil.RespondWithError(w, http.StatusInternalServerError, "")
		return
	}

	token.GetUserId()

	h.carService.CreateCar(r.Context(), token.GetUserId(), car.Car{
		Make:  req.Make,
		Model: req.Model,
		Year:  req.Year,
		Trim:  req.Trim,
		VIN:   req.VIN,
	})

	httputil.RespondWithJSON(w, http.StatusOK, req)
}
