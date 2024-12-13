package ports

import (
	"context"
	"database/sql"

	"github.com/hilmiikhsan/multifinance-service/internal/module/transaction/dto"
	"github.com/hilmiikhsan/multifinance-service/internal/module/transaction/entity"
)

type TransactionRepository interface {
	InsertNewTransaction(ctx context.Context, tx *sql.Tx, data *entity.Transaction) error
}

type TransactionService interface {
	CreateTransaction(ctx context.Context, req *dto.CreateTransactionRequest) error
}
