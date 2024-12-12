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
	"github.com/hilmiikhsan/multifinance-service/pkg/utils"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var _ authPorts.AuthService = &authService{}

type authService struct {
	db                 *sqlx.DB
	customerRepository customerPorts.CustomerRepository
	redisDB            redisPorts.RedisRepository
}

func NewUserService(db *sqlx.DB, customerRepository customerPorts.CustomerRepository, redisDB redisPorts.RedisRepository) *authService {
	return &authService{
		db:                 db,
		customerRepository: customerRepository,
		redisDB:            redisDB,
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
