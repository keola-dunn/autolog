package utils

import "github.com/golang-jwt/jwt/v5"

type AutologAPIJWTClaims struct {
	jwt.RegisteredClaims
}
