package cars

import (
	"net/http"

	autologjwt "github.com/keola-dunn/autolog/cmd/autolog-api/jwt"
	"github.com/keola-dunn/autolog/internal/httputil"
	"github.com/keola-dunn/autolog/internal/logger"
)

//type createCarRequest struct{}

//type createCarResponse struct {}

func (h *CarsHandler) CreateCar(w http.ResponseWriter, r *http.Request) {
	logEntry := logger.GetLogEntry(r)

	authHeader := r.Header.Get("Authorization")
	ok, token, err := autologjwt.VerifyToken(authHeader, "")
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

	userId := token.GetUserId()
	httputil.RespondWithJSON(w, http.StatusOK, userId)
}
