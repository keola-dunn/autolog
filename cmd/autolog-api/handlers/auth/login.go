package auth

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/keola-dunn/autolog/internal/httputil"
)

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, pass, ok := r.BasicAuth()
	if !ok {
		httputil.RespondWithError(w, http.StatusUnauthorized, "missing required user/pass")
		return
	}

	valid, err := h.userService.ValidCredentials(ctx, user, pass)
	if err != nil {
		httputil.RespondWithError(w, http.StatusInternalServerError, "")
		return
	}
	if !valid {
		httputil.RespondWithError(w, http.StatusUnauthorized, "")
	}

	claims := jwt.RegisteredClaims{
		Issuer:    "",
		Subject:   "",
		Audience:  jwt.ClaimStrings{},
		ExpiresAt: jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ID:        "",
	}

	type myCustomClaims struct {
		CustomClaim string
		jwt.RegisteredClaims
	}

	myClaims := myCustomClaims{
		CustomClaim:      "",
		RegisteredClaims: claims,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, myClaims)

	jwtToken, err := token.SignedString(h.jwtSecret)
	if err != nil {

	}

	w.Write([]byte(jwtToken))
}
