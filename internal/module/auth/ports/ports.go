package ports

import (
	"context"

	"github.com/hilmiikhsan/multifinance-service/internal/module/auth/dto"
)

type AuthService interface {
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error)
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error)
	RefreshToken(ctx context.Context, accessToken string) (*dto.RefreshTokenResponse, error)
}
