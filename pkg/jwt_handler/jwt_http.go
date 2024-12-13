package jwt_handler

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/hilmiikhsan/multifinance-service/constants"
	"github.com/hilmiikhsan/multifinance-service/internal/infrastructure/config"
	redisPorts "github.com/hilmiikhsan/multifinance-service/internal/infrastructure/redis"
	"github.com/hilmiikhsan/multifinance-service/pkg/err_msg"
	"github.com/rs/zerolog/log"
)

var _ JWT = &jwtHandler{}

type jwtHandler struct {
	db redisPorts.RedisRepository
}

func NewJWT(db redisPorts.RedisRepository) *jwtHandler {
	return &jwtHandler{
		db: db,
	}
}

func (j *jwtHandler) GenerateTokenString(ctx context.Context, payload CostumClaimsPayload) (string, error) {
	MapTypeToken := make(map[string]time.Duration)

	tokenDuration, err := time.ParseDuration(config.Envs.Guard.JwtTokenExpiration)
	if err != nil {
		log.Error().Err(err).Msg("jwthandler::GenerateTokenString - Error while parsing token duration")
		return "", err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	MapTypeToken[constants.AccessTokenType] = tokenDuration

	refreshTokenDuration, err := time.ParseDuration(config.Envs.Guard.JwtRefreshTokenExpiration)
	if err != nil {
		log.Error().Err(err).Msg("jwthandler::GenerateTokenString - Error while parsing refresh token duration")
		return "", err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	MapTypeToken[constants.RefreshTokenType] = refreshTokenDuration

	expireTime := time.Now().Add(MapTypeToken[payload.TokenType])

	claims := CustomClaims{
		UserId:   payload.UserId,
		Email:    payload.Email,
		FullName: payload.FullName,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "user",
			Issuer:    config.Envs.App.Name,
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			NotBefore: jwt.NewNumericDate(time.Now().UTC()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	tokenString, err := token.SignedString([]byte(config.Envs.Guard.JwtPrivateKey))
	if err != nil {
		log.Error().Err(err).Msg("jwthandler::GenerateTokenString - Error while signing token")
		return "", err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	expirationDuration := time.Until(expireTime)

	key := fmt.Sprintf("%s:%s", payload.Nik, payload.TokenType)

	err = j.db.Set(ctx, key, claims.UserId, expirationDuration)
	if err != nil {
		log.Error().Err(err).Msg("jwthandler::GenerateTokenString - Error while saving token to Redis")
		return "", err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return tokenString, nil
}

func (j *jwtHandler) ParseTokenString(ctx context.Context, tokenString, username, tokenType string) (*CustomClaims, error) {
	claims := &CustomClaims{}

	key := fmt.Sprintf("%s:%s", username, tokenType)

	_, err := j.db.Get(ctx, key)
	if err != nil {
		log.Error().Err(err).Msg("jwthandler::ParseTokenString - Token not found in Redis")
		return nil, err_msg.NewCustomErrors(fiber.StatusUnauthorized, err_msg.WithMessage(constants.ErrTokenAlreadyExpired))
	}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Envs.Guard.JwtPrivateKey), nil
	})
	if err != nil {
		log.Error().Err(err).Msg("jwthandler::ParseTokenString - Error while parsing token")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	if !token.Valid {
		log.Error().Msg("jwthandler::ParseTokenString - Invalid token")
		return nil, err_msg.NewCustomErrors(fiber.StatusUnauthorized, err_msg.WithMessage(constants.ErrTokenAlreadyExpired))
	}

	return claims, nil
}
