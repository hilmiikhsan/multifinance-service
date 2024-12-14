package ports

import (
	"context"
	"database/sql"

	"github.com/hilmiikhsan/multifinance-service/internal/module/transaction/dto"
	"github.com/hilmiikhsan/multifinance-service/internal/module/transaction/entity"
)

//go:generate mockgen -source=ports.go -destination=../service/service_mock_test.go -package=service
type TransactionRepository interface {
	InsertNewTransaction(ctx context.Context, tx *sql.Tx, data *entity.Transaction) error
	FindTransactionByIdAndCustomerID(ctx context.Context, id, customerID int) (*entity.Transaction, error)
	FindTransactionByCustomerID(ctx context.Context, req *dto.GetHistoryListTransactionRequest, customerID int) (*dto.GetHistoryListTransactionResponse, error)
}

//go:generate mockgen -source=ports.go -destination=../handler/rest/handler_mock_test.go -package=rest
type TransactionService interface {
	CreateTransaction(ctx context.Context, req *dto.CreateTransactionRequest) error
	GetDetailTransaction(ctx context.Context, id, customerID int) (*dto.GetDetailTransactionResponse, error)
	GetHistoryListTransction(ctx context.Context, req *dto.GetHistoryListTransactionRequest, customerID int) (*dto.GetHistoryListTransactionResponse, error)
}
