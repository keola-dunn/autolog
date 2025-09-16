package jwt

import "context"

type contextKey string

const (
	contextKeyClaims = contextKey("claims")
)

func SetClaimsInContext(ctx context.Context, claims AutologAPIJWTClaims) context.Context {
	return context.WithValue(ctx, contextKeyClaims, claims)
}

func GetClaimsFromContext(ctx context.Context) (AutologAPIJWTClaims, bool) {
	if c := ctx.Value(contextKeyClaims); c != nil {
		if claims, ok := c.(AutologAPIJWTClaims); ok {
			return claims, true
		}
	}

	return AutologAPIJWTClaims{}, false

}
