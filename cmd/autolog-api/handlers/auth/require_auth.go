package auth

import (
	"net/http"
	"strings"

	autologjwt "github.com/keola-dunn/autolog/cmd/autolog-api/jwt"
	"github.com/keola-dunn/autolog/internal/httputil"
)

// RequireAuthentication is a middleware that requires the request to be authenticated
func (a *AuthHandler) RequireTokenAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if strings.TrimSpace(authHeader) == "" {
			a.logger.Warn("missing authentication header",
				"referer", r.Header.Get("referer"),
				"user-agent", r.Header.Get("user-agent"),
				"x-forwarded-for", r.Header.Get("X-Forwarded-For"))
			httputil.RespondWithError(w, http.StatusUnauthorized, "")
			return
		}

		splitToken := strings.Split(authHeader, "Bearer ")
		if len(splitToken) != 2 || !strings.Contains(authHeader, "Bearer") {
			a.logger.Warn("invalid authentication header",
				"header", authHeader,
				"referer", r.Header.Get("referer"),
				"user-agent", r.Header.Get("user-agent"),
				"x-forwarded-for", r.Header.Get("X-Forwarded-For"))
			httputil.RespondWithError(w, http.StatusUnauthorized, "")
			return
		}

		token := splitToken[1]

		valid, _, err := autologjwt.VerifyToken(token, a.jwtSecret)
		if err != nil {
			a.logger.Error("failed to verify token", err)
			httputil.RespondWithError(w, http.StatusInternalServerError, "")
			return
		}
		if !valid {
			// log failed auth attempts
			a.logger.Warn("invalid token provided",
				"token", token,
				"referer", r.Header.Get("referer"),
				"user-agent", r.Header.Get("user-agent"),
				"x-forwarded-for", r.Header.Get("X-Forwarded-For"))

			httputil.RespondWithError(w, http.StatusUnauthorized, "")
			return
		}

		next.ServeHTTP(w, r)
	})
}
