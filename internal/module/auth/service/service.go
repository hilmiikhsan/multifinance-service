package service

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/hilmiikhsan/multifinance-service/constants"
	redisPorts "github.com/hilmiikhsan/multifinance-service/internal/infrastructure/redis"
	"github.com/hilmiikhsan/multifinance-service/internal/module/auth/dto"
	authPorts "github.com/hilmiikhsan/multifinance-service/internal/module/auth/ports"
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
	db                 *sqlx.DB
	customerRepository customerPorts.CustomerRepository
	redisDB            redisPorts.RedisRepository
	jwt                jwt_handler.JWT
}

func NewUserService(db *sqlx.DB, customerRepository customerPorts.CustomerRepository, redisDB redisPorts.RedisRepository, jwt jwt_handler.JWT) *authService {
	return &authService{
		db:                 db,
		customerRepository: customerRepository,
		redisDB:            redisDB,
		jwt:                jwt,
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
	salary, _ := strconv.ParseFloat(req.Salary, 64)

	result, err := s.customerRepository.InsertNewUser(ctx, &entity.Customer{
		Nik:             req.Nik,
		Email:           req.Email,
		Password:        req.Password,
		FullName:        req.FullName,
		LegalName:       req.LegalName,
		BirthPlace:      req.BirthPlace,
		BirthDate:       birthDate,
		Salary:          salary,
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
		UserId:    customerData.ID,
		Nik:       customerData.Nik,
		Email:     customerData.Email,
		FullName:  customerData.FullName,
		TokenType: constants.AccessTokenType,
	})
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::Login - Failed to generate token string")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	refreshToken, err := s.jwt.GenerateTokenString(ctx, jwt_handler.CostumClaimsPayload{
		UserId:    customerData.ID,
		Nik:       customerData.Nik,
		Email:     customerData.Email,
		FullName:  customerData.FullName,
		TokenType: constants.RefreshTokenType,
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
		UserId:    id,
		Nik:       claims.Nik,
		Email:     claims.Email,
		FullName:  claims.FullName,
		TokenType: constants.AccessTokenType,
	})
	if err != nil {
		log.Error().Err(err).Any("payload", claims).Msg("service::RefreshToken - Failed to generate token string")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	res.Token = token

	return res, nil
}
