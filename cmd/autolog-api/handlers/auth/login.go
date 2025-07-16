package auth

import (
	"net/http"

	"github.com/keola-dunn/autolog/internal/httputil"
)

type LoginResponse struct {
	JWT string `json:"jwt"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, pass, ok := r.BasicAuth()
	if !ok {
		httputil.RespondWithError(w, http.StatusUnauthorized, "missing required user/pass")
		return
	}

	valid, userId, err := h.userService.ValidateCredentials(ctx, user, pass)
	if err != nil {
		h.logger.Error("failed to validate credentials", err)
		httputil.RespondWithError(w, http.StatusInternalServerError, "")
		return
	}
	if !valid {
		httputil.RespondWithError(w, http.StatusUnauthorized, "")
		return
	}

	jwtToken, err := h.createJWT(userId)
	if err != nil {
		h.logger.Error("failed to create jwt", err)
		httputil.RespondWithError(w, http.StatusInternalServerError, "")
		return
	}

	httputil.RespondWithJSON(w, http.StatusOK, LoginResponse{
		JWT: jwtToken,
	})
}
