package service

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/hilmiikhsan/multifinance-service/constants"
	dtoLimit "github.com/hilmiikhsan/multifinance-service/internal/module/credit_limit/dto"
	"github.com/hilmiikhsan/multifinance-service/internal/module/customer/dto"
	customerPorts "github.com/hilmiikhsan/multifinance-service/internal/module/customer/ports"
	"github.com/hilmiikhsan/multifinance-service/pkg/err_msg"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var _ customerPorts.CustomerService = &customerService{}

type customerService struct {
	db                 *sqlx.DB
	customerRepository customerPorts.CustomerRepository
}

func NewCustomerService(db *sqlx.DB, customerRepository customerPorts.CustomerRepository) *customerService {
	return &customerService{
		db:                 db,
		customerRepository: customerRepository,
	}
}

func (s *customerService) GetCustomerProfile(ctx context.Context, id int) (*dto.GetCustomerProfileResponse, error) {
	customer, err := s.customerRepository.FindCustomerByID(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrUserNotFound) {
			log.Error().Err(err).Msg("service::GetCustomerProfile - Customer not found")
			return nil, err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrUserNotFound))
		}

		log.Error().Err(err).Msg("service::GetCustomerProfile - Failed to find customer by ID")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	var limits []dtoLimit.CreditLimit
	for _, limit := range customer.Limits {
		limits = append(limits, dtoLimit.CreditLimit{
			Tenor:       limit.TenorMonth,
			LimitAmount: limit.LimitAmount,
		})
	}

	return &dto.GetCustomerProfileResponse{
		ID:              customer.ID,
		Nik:             customer.Nik,
		FullName:        customer.FullName,
		LegalName:       customer.LegalName,
		BirthPlace:      customer.BirthPlace,
		BirthDate:       customer.BirthDate.Format(constants.DateTimeFormat),
		Salary:          customer.Salary,
		KtpPhotoPath:    customer.KtpPhotoPath,
		SelfiePhotoPath: customer.SelfiePhotoPath,
		Limits:          limits,
		CreatedAt:       customer.CreatedAt.Format(constants.DateTimeFormat),
		UpdatedAt:       customer.UpdatedAt.Format(constants.DateTimeFormat),
	}, nil
}
