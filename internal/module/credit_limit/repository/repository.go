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
