package service

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/hilmiikhsan/multifinance-service/constants"
	creditLimitPorts "github.com/hilmiikhsan/multifinance-service/internal/module/credit_limit/ports"
	"github.com/hilmiikhsan/multifinance-service/internal/module/transaction/dto"
	"github.com/hilmiikhsan/multifinance-service/internal/module/transaction/entity"
	transactionPorts "github.com/hilmiikhsan/multifinance-service/internal/module/transaction/ports"
	"github.com/hilmiikhsan/multifinance-service/pkg/err_msg"
	"github.com/hilmiikhsan/multifinance-service/pkg/utils"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var _ transactionPorts.TransactionService = &transactionService{}

type transactionService struct {
	db                    *sqlx.DB
	transactionRepository transactionPorts.TransactionRepository
	creditLimitRepository creditLimitPorts.CreditLimitRepository
}

func NewTransactionService(db *sqlx.DB, transactionRepository transactionPorts.TransactionRepository, creditLimitRepository creditLimitPorts.CreditLimitRepository) *transactionService {
	return &transactionService{
		db:                    db,
		transactionRepository: transactionRepository,
		creditLimitRepository: creditLimitRepository,
	}
}

func (s *transactionService) CreateTransaction(ctx context.Context, req *dto.CreateTransactionRequest) error {
	// Step 0: Begin transaction
	tx, err := s.db.Begin()
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::CreateTransaction - Failed to begin transaction")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Error().Err(rollbackErr).Any("payload", req).Msg("service::CreateTransaction - Failed to rollback transaction")
			}
		}
	}()

	// Step 1: Validate tenor and credit limit with locking
	creditLimit, err := s.creditLimitRepository.FindLimitByCustomerAndTenor(ctx, tx, req.CustomerID, req.TenorMonth)
	if err != nil {
		log.Error().Err(err).Int("customer_id", req.CustomerID).Msg("service::CreateTransaction - Failed to find credit limit for customer and tenor")
		return err_msg.NewCustomErrors(fiber.StatusBadRequest, err_msg.WithMessage(constants.ErrInvalidOrCreditLimit))
	}

	// Step 2: Validate OnTheRoadPrice does not exceed limit amount
	if req.OnTheRoadPrice > int(creditLimit.LimitAmount) {
		log.Warn().
			Int("customer_id", req.CustomerID).
			Int("on_the_road_price", req.OnTheRoadPrice).
			Float64("limit_amount", creditLimit.LimitAmount).
			Msg("service::CreateTransaction - On the road price exceeds credit limit")
		return err_msg.NewCustomErrors(fiber.StatusBadRequest, err_msg.WithMessage(constants.ErrOnTheRoadPriceExceedLimit))
	}

	// Step 3: Generate contract number
	contractNumber := utils.GenerateContractNumber(req.CustomerID)

	// Step 4: Calculate fees and amounts
	adminFee := utils.CalculateAdminFee(req.OnTheRoadPrice)
	interestAmount := utils.CalculateInterest(req.OnTheRoadPrice, req.TenorMonth)
	installmentAmount := utils.CalculateInstallment(req.OnTheRoadPrice, interestAmount, req.TenorMonth)

	// Step 5: Create transaction entity
	transaction := &entity.Transaction{
		CustomerID:        req.CustomerID,
		ContractNumber:    contractNumber,
		OnTheRoadPrice:    float64(req.OnTheRoadPrice),
		AdminFee:          float64(adminFee),
		InstallmentAmount: float64(installmentAmount),
		InterestAmount:    float64(interestAmount),
		AssetName:         req.AssetName,
	}

	// Step 6: Insert transaction into database
	err = s.transactionRepository.InsertNewTransaction(ctx, tx, transaction)
	if err != nil {
		log.Error().Err(err).Msg("service::CreateTransaction - Failed to insert new transaction")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// Step 7: Commit transaction
	err = tx.Commit()
	if err != nil {
		log.Error().Err(err).Msg("service::CreateTransaction - Failed to commit transaction")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	log.Info().Str("contract_number", contractNumber).Msg("service::CreateTransaction - Transaction created successfully")
	return nil
}
