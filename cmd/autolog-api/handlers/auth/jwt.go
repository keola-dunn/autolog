package auth

import (
	"fmt"
	"time"

	autologjwt "github.com/keola-dunn/autolog/cmd/autolog-api/jwt"
)

func (h *AuthHandler) createJWT(userId string) (string, error) {
	now := h.calendarService.NowUTC()

	tokenId, err := h.randomGenerator.RandomUUID()
	if err != nil {
		return "", fmt.Errorf("failed to create random token id: %w", err)
	}

	jwtToken, err := autologjwt.CreateJWT(autologjwt.CreateJWTInput{
		Issuer:      h.jwtIssuer,
		UserId:      userId,
		IssuedAt:    now,
		ExpiresAt:   now.Add(time.Duration(h.jwtExpiryLengthMinutes) * time.Minute),
		NotBefore:   now,
		Id:          tokenId,
		TokenSecret: h.jwtSecret,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create jwt: %w", err)
	}

	return jwtToken, nil
}
