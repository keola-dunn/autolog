package auth

import (
	"net/http"

	"github.com/keola-dunn/autolog/internal/httputil"
	"github.com/keola-dunn/autolog/internal/jwt"
	"github.com/keola-dunn/autolog/internal/logger"
)

func (a *AuthHandler) GetWellKnownJWKS(w http.ResponseWriter, r *http.Request) {
	logEntry := logger.GetLogEntry(r)

	jwk, err := jwt.ConvertPublicKeyPEMToJWK("autolog-public-key", a.jwtPublicKey)
	if err != nil {
		logEntry.Error("failed to convert public key pem to jwk", err)
		httputil.RespondWithError(w, http.StatusInternalServerError, "")
		return
	}

	httputil.RespondWithJSON(w, http.StatusOK, jwt.JWKS{
		Keys: []jwt.JWK{
			jwk,
		},
	})
}
