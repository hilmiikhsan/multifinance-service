package repository

import (
	"context"
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/hilmiikhsan/multifinance-service/constants"
	"github.com/hilmiikhsan/multifinance-service/internal/module/credit_limit/entity"
	"github.com/hilmiikhsan/multifinance-service/internal/module/credit_limit/ports"
	"github.com/hilmiikhsan/multifinance-service/pkg/err_msg"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var _ ports.CreditLimitRepository = &creditLimitRepository{}

type creditLimitRepository struct {
	db *sqlx.DB
}

func NewCreditLimitRepository(db *sqlx.DB) *creditLimitRepository {
	return &creditLimitRepository{
		db: db,
	}
}

func (r *creditLimitRepository) InsertNewCreditLimit(ctx context.Context, tx *sql.Tx, data *entity.CreditLimit) error {
	_, err := tx.ExecContext(ctx, r.db.Rebind(queryInsertNewCreditLimit),
		data.CustomerID,
		data.TenorMonth,
		data.LimitAmount,
	)
	if err != nil {
		log.Error().Err(err).Any("payload", data).Msg("repository::InsertNewCreditLimit - Failed to insert new credit limit")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return nil
}

func (r *creditLimitRepository) FindCreditLimitByCustomerID(ctx context.Context, customerID int) (*[]entity.Limits, error) {
	var limits []entity.Limits

	err := r.db.SelectContext(ctx, &limits, r.db.Rebind(queryFindCreditLimitByCustomerID), customerID)
	if err != nil {
		log.Error().Err(err).Int("customerID", customerID).Msg("repository::FindCreditLimitByCustomerID - Failed to find credit limit by customer ID")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return &limits, nil
}

func (r *creditLimitRepository) FindLimitByCustomerAndTenor(ctx context.Context, tx *sql.Tx, customerID int, tenorMonth int) (*entity.Limits, error) {
	var limit entity.Limits

	err := tx.QueryRowContext(ctx, queryLockCreditLimitByCustomerAndTenor, customerID, tenorMonth).Scan(&limit.TenorMonth, &limit.LimitAmount)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Error().
				Err(err).
				Int("customer_id", customerID).
				Int("tenor_month", tenorMonth).
				Msg("repository::FindLimitByCustomerAndTenor - No credit limit found")
			return nil, err_msg.NewCustomErrors(fiber.StatusBadRequest, err_msg.WithMessage("Invalid tenor or customer ID"))
		}
		log.Error().
			Err(err).
			Int("customer_id", customerID).
			Int("tenor_month", tenorMonth).
			Msg("repository::FindLimitByCustomerAndTenor - Failed to fetch credit limit")
		return nil, err
	}

	return &limit, nil
}
