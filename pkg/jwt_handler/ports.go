package jwt_handler

import "context"

type JWT interface {
	GenerateTokenString(ctx context.Context, payload CostumClaimsPayload) (string, error)
	ParseTokenString(ctx context.Context, tokenString, username, tokenType string) (*CustomClaims, error)
}
