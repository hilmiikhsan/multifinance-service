package repository

import (
	"context"
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/hilmiikhsan/multifinance-service/constants"
	"github.com/hilmiikhsan/multifinance-service/internal/module/transaction/entity"
	"github.com/hilmiikhsan/multifinance-service/internal/module/transaction/ports"
	"github.com/hilmiikhsan/multifinance-service/pkg/err_msg"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var _ ports.TransactionRepository = &transactionRepository{}

type transactionRepository struct {
	db *sqlx.DB
}

func NewTransactionRepository(db *sqlx.DB) *transactionRepository {
	return &transactionRepository{
		db: db,
	}
}

func (r *transactionRepository) InsertNewTransaction(ctx context.Context, tx *sql.Tx, data *entity.Transaction) error {
	_, err := tx.ExecContext(ctx, r.db.Rebind(queryInsertNewTransaction),
		data.CustomerID,
		data.ContractNumber,
		data.OnTheRoadPrice,
		data.AdminFee,
		data.InstallmentAmount,
		data.InterestAmount,
		data.AssetName,
	)
	if err != nil {
		log.Error().Err(err).Msg("repository::CreateTransaction - Failed to insert new transaction")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return nil
}
