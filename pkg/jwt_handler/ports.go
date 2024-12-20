package jwt_handler

import "context"

//go:generate mockgen -source=ports.go -destination=service_jwt_mock_test.go -package=jwt_handler
type JWT interface {
	GenerateTokenString(ctx context.Context, payload CostumClaimsPayload) (string, error)
	ParseTokenString(ctx context.Context, tokenString string) (*CustomClaims, error)
}
