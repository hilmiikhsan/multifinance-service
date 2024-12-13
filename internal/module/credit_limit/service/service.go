package service

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/hilmiikhsan/multifinance-service/constants"
	"github.com/hilmiikhsan/multifinance-service/internal/module/credit_limit/dto"
	creditLimitPorts "github.com/hilmiikhsan/multifinance-service/internal/module/credit_limit/ports"
	"github.com/hilmiikhsan/multifinance-service/pkg/err_msg"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var _ creditLimitPorts.CreditLimitService = &creditLimitService{}

type creditLimitService struct {
	db                    *sqlx.DB
	creditLimitRepository creditLimitPorts.CreditLimitRepository
}

func NewCreditLimitService(db *sqlx.DB, creditLimitRepository creditLimitPorts.CreditLimitRepository) *creditLimitService {
	return &creditLimitService{
		db:                    db,
		creditLimitRepository: creditLimitRepository,
	}
}

func (s *creditLimitService) GetCreditLimits(ctx context.Context, customerID int) (*[]dto.GetCreditLimitsResponse, error) {
	res := new([]dto.GetCreditLimitsResponse)

	creditLimits, err := s.creditLimitRepository.FindCreditLimitByCustomerID(ctx, customerID)
	if err != nil {
		log.Error().Err(err).Int("customerID", customerID).Msg("service::GetCreditLimits - Failed to find credit limits by customer ID")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	for _, creditLimit := range *creditLimits {
		*res = append(*res, dto.GetCreditLimitsResponse{
			Tenor:       creditLimit.TenorMonth,
			LimitAmount: creditLimit.LimitAmount,
		})
	}

	return res, nil
}
