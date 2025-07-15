package auth

import "github.com/keola-dunn/autolog/internal/service/user"

type AuthHandler struct {
	jwtSecret   string
	userService user.ServiceIface
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}
