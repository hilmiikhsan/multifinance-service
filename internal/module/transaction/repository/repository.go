package repository

import (
	"context"
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/hilmiikhsan/multifinance-service/constants"
	"github.com/hilmiikhsan/multifinance-service/internal/module/transaction/dto"
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

func (r *transactionRepository) FindTransactionByIdAndCustomerID(ctx context.Context, id, customerID int) (*entity.Transaction, error) {
	var (
		res = new(entity.Transaction)
	)

	err := r.db.GetContext(ctx, res, r.db.Rebind(queryFindTransactionByIdAndCustomerID), id, customerID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Error().Err(err).Msg("repository::FindTransactionByIdAndCustomerID - Transaction not found")
			return nil, err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrTransactionNotFound))
		}

		log.Error().Err(err).Msg("repository::FindTransactionByIdAndCustomerID - Failed to find transaction by ID and customer ID")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return res, nil
}

func (r *transactionRepository) FindTransactionByCustomerID(ctx context.Context, req *dto.GetHistoryListTransactionRequest, customerID int) (*dto.GetHistoryListTransactionResponse, error) {
	var (
		resp       = new(dto.GetHistoryListTransactionResponse)
		data       = make([]dto.HistoryListTransactionItem, 0, req.Paginate)
		query      = queryFindTransactionByCustomerID
		countQuery = queryCountTransactionByCustomerID
	)

	var totalData int
	countQuery, countArgs, err := sqlx.Named(countQuery, map[string]interface{}{
		"customer_id": customerID,
	})
	if err != nil {
		log.Error().Err(err).Msg("repository::FindTransactionByCustomerID - Failed to bind named query for count")
		return nil, err
	}

	countQuery = r.db.Rebind(countQuery)
	err = r.db.GetContext(ctx, &totalData, countQuery, countArgs...)
	if err != nil {
		log.Error().Err(err).Msg("repository::FindTransactionByCustomerID - Failed to count transactions")
		return nil, err
	}

	query, args, err := sqlx.Named(query, map[string]interface{}{
		"customer_id": customerID,
		"limit":       req.Paginate,
		"offset":      req.Paginate * (req.Page - 1),
	})
	if err != nil {
		log.Error().Err(err).Msg("repository::FindTransactionByCustomerID - Failed to bind named query")
		return nil, err
	}

	query = r.db.Rebind(query)
	err = r.db.SelectContext(ctx, &data, query, args...)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repository::FindTransactionByCustomerID - Failed to find transactions")
		return nil, err
	}

	resp.Items = data
	resp.Meta.TotalData = totalData
	resp.Meta.CountTotalPage(req.Page, req.Paginate, totalData)

	return resp, nil
}
