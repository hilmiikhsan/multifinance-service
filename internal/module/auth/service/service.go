package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/hilmiikhsan/multifinance-service/constants"
	redisPorts "github.com/hilmiikhsan/multifinance-service/internal/infrastructure/redis"
	"github.com/hilmiikhsan/multifinance-service/internal/middleware"
	"github.com/hilmiikhsan/multifinance-service/internal/module/auth/dto"
	authPorts "github.com/hilmiikhsan/multifinance-service/internal/module/auth/ports"
	creditLimitEntity "github.com/hilmiikhsan/multifinance-service/internal/module/credit_limit/entity"
	creditLimitPorts "github.com/hilmiikhsan/multifinance-service/internal/module/credit_limit/ports"
	"github.com/hilmiikhsan/multifinance-service/internal/module/customer/entity"
	customerPorts "github.com/hilmiikhsan/multifinance-service/internal/module/customer/ports"
	"github.com/hilmiikhsan/multifinance-service/pkg/err_msg"
	"github.com/hilmiikhsan/multifinance-service/pkg/jwt_handler"
	"github.com/hilmiikhsan/multifinance-service/pkg/utils"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var _ authPorts.AuthService = &authService{}

type authService struct {
	db                    *sqlx.DB
	customerRepository    customerPorts.CustomerRepository
	redisDB               redisPorts.RedisRepository
	jwt                   jwt_handler.JWT
	creditLimitRepository creditLimitPorts.CreditLimitRepository
}

func NewUserService(db *sqlx.DB, customerRepository customerPorts.CustomerRepository, redisDB redisPorts.RedisRepository, jwt jwt_handler.JWT, creditLimitRepository creditLimitPorts.CreditLimitRepository) *authService {
	return &authService{
		db:                    db,
		customerRepository:    customerRepository,
		redisDB:               redisDB,
		jwt:                   jwt,
		creditLimitRepository: creditLimitRepository,
	}
}

func (s *authService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error) {
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::Register - Failed to hash password")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	req.Password = hashedPassword

	birthDate, _ := time.Parse(constants.DateTimeFormat, req.BirthDate)

	tx, err := s.db.Begin()
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::Register - Failed to begin transaction")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Error().Err(rollbackErr).Any("payload", req).Msg("service::Register - Failed to rollback transaction")
			}
		}
	}()

	result, err := s.customerRepository.InsertNewUser(ctx, tx, &entity.Customer{
		Nik:             req.Nik,
		Email:           req.Email,
		Password:        req.Password,
		FullName:        req.FullName,
		LegalName:       req.LegalName,
		BirthPlace:      req.BirthPlace,
		BirthDate:       birthDate,
		Salary:          float64(req.Salary),
		KtpPhotoPath:    req.KtpPhotoPath,
		SelfiePhotoPath: req.SelfiePhotoPath,
	})
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrNikAlreadyRegistered) {
			log.Error().Err(err).Any("payload", req).Msg("service::Register - Failed to insert new user")
			return nil, err_msg.NewCustomErrors(fiber.StatusConflict, err_msg.WithMessage(constants.ErrNikAlreadyRegistered))
		}

		if strings.Contains(err.Error(), constants.ErrEmailAlreadyRegistered) {
			log.Error().Err(err).Any("payload", req).Msg("service::Register - Failed to insert new user")
			return nil, err_msg.NewCustomErrors(fiber.StatusConflict, err_msg.WithMessage(constants.ErrEmailAlreadyRegistered))
		}

		log.Error().Err(err).Any("payload", req).Msg("service::Register - Failed to insert new user")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	var defaultLimits []creditLimitEntity.CreditLimit
	switch {
	case req.Salary < 5000000:
		defaultLimits = []creditLimitEntity.CreditLimit{
			{CustomerID: result.ID, TenorMonth: 1, LimitAmount: 100000.00},
			{CustomerID: result.ID, TenorMonth: 2, LimitAmount: 200000.00},
			{CustomerID: result.ID, TenorMonth: 3, LimitAmount: 500000.00},
			{CustomerID: result.ID, TenorMonth: 6, LimitAmount: 700000.00},
		}
	case req.Salary <= 10000000:
		defaultLimits = []creditLimitEntity.CreditLimit{
			{CustomerID: result.ID, TenorMonth: 1, LimitAmount: 200000.00},
			{CustomerID: result.ID, TenorMonth: 2, LimitAmount: 400000.00},
			{CustomerID: result.ID, TenorMonth: 3, LimitAmount: 800000.00},
			{CustomerID: result.ID, TenorMonth: 6, LimitAmount: 1200000.00},
		}
	default: // Salary > 10 juta
		defaultLimits = []creditLimitEntity.CreditLimit{
			{CustomerID: result.ID, TenorMonth: 1, LimitAmount: 500000.00},
			{CustomerID: result.ID, TenorMonth: 2, LimitAmount: 1000000.00},
			{CustomerID: result.ID, TenorMonth: 3, LimitAmount: 1500000.00},
			{CustomerID: result.ID, TenorMonth: 6, LimitAmount: 2000000.00},
		}
	}

	for _, limit := range defaultLimits {
		if err := s.creditLimitRepository.InsertNewCreditLimit(ctx, tx, &limit); err != nil {
			log.Error().Err(err).Any("payload", limit).Msg("service::Register - Failed to insert new credit limit")
			return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
		}
	}

	if err := tx.Commit(); err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::Register - Failed to commit transaction")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return &dto.RegisterResponse{
		ID:    result.ID,
		Email: result.Email,
	}, nil
}

func (s *authService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	var (
		res = new(dto.LoginResponse)
	)

	customerData, err := s.customerRepository.FindCustomerByEmail(ctx, req.Email)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::Login - Failed to find user")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	if customerData == nil {
		log.Error().Any("payload", req).Msg("service::Login - Email not found")
		return nil, err_msg.NewCustomErrors(fiber.StatusUnprocessableEntity, err_msg.WithMessage(constants.ErrEmailOrPasswordIsIncorrect))
	}

	if !utils.ComparePassword(customerData.Password, req.Password) {
		log.Error().Any("payload", req).Msg("service::Login - Password is incorrect")
		return nil, err_msg.NewCustomErrors(fiber.StatusUnprocessableEntity, err_msg.WithMessage(constants.ErrEmailOrPasswordIsIncorrect))
	}

	token, err := s.jwt.GenerateTokenString(ctx, jwt_handler.CostumClaimsPayload{
		CustomerID: customerData.ID,
		Nik:        customerData.Nik,
		Email:      customerData.Email,
		FullName:   customerData.FullName,
		TokenType:  constants.AccessTokenType,
	})
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::Login - Failed to generate token string")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	refreshToken, err := s.jwt.GenerateTokenString(ctx, jwt_handler.CostumClaimsPayload{
		CustomerID: customerData.ID,
		Nik:        customerData.Nik,
		Email:      customerData.Email,
		FullName:   customerData.FullName,
		TokenType:  constants.RefreshTokenType,
	})
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::Login - Failed to generate token string")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	res.ID = customerData.ID
	res.Email = customerData.Email
	res.FullName = customerData.FullName
	res.Token = token
	res.RefreshToken = refreshToken

	return res, nil
}

func (s *authService) RefreshToken(ctx context.Context, accessToken string) (*dto.RefreshTokenResponse, error) {
	var (
		res = new(dto.RefreshTokenResponse)
	)

	claims, err := s.jwt.ParseTokenString(ctx, accessToken)
	if err != nil {
		log.Error().Err(err).Any("access_token", accessToken).Msg("service::RefreshToken - Failed to parse access token")
		return nil, err_msg.NewCustomErrors(fiber.StatusUnauthorized, err_msg.WithMessage(constants.ErrInvalidAccessToken))
	}

	id, _ := strconv.Atoi(claims.ID)

	token, err := s.jwt.GenerateTokenString(ctx, jwt_handler.CostumClaimsPayload{
		CustomerID: int64(id),
		Nik:        claims.Nik,
		Email:      claims.Email,
		FullName:   claims.FullName,
		TokenType:  constants.AccessTokenType,
	})
	if err != nil {
		log.Error().Err(err).Any("payload", claims).Msg("service::RefreshToken - Failed to generate token string")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	res.Token = token

	return res, nil
}

func (s *authService) Logout(ctx context.Context, accessToken string, locals *middleware.Locals) error {
	key := fmt.Sprintf("%s:%s", locals.Nik, constants.AccessTokenType)

	_, err := s.redisDB.Get(ctx, key)
	if err != nil {
		log.Error().Err(err).Msg("jwthandler::ParseTokenString - Token not found in Redis")
		return err_msg.NewCustomErrors(fiber.StatusUnauthorized, err_msg.WithMessage(constants.ErrTokenAlreadyExpired))
	}

	claims, err := s.jwt.ParseTokenString(ctx, accessToken)
	if err != nil {
		log.Error().Err(err).Any("access_token", accessToken).Msg("service::Logout - Failed to parse access token")
		return err_msg.NewCustomErrors(fiber.StatusUnauthorized, err_msg.WithMessage(constants.ErrInvalidAccessToken))
	}

	key = fmt.Sprintf("%s:%s", claims.Nik, constants.AccessTokenType)

	err = s.redisDB.Del(ctx, key)
	if err != nil {
		log.Error().Err(err).Any("access_token", accessToken).Msg("service::Logout - Failed to set access token to redis")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return nil
}
