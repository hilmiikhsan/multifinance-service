package ports

import (
	"context"
	"database/sql"

	"github.com/hilmiikhsan/multifinance-service/internal/module/credit_limit/dto"
	"github.com/hilmiikhsan/multifinance-service/internal/module/credit_limit/entity"
)

//go:generate mockgen -source=ports.go -destination=../service/service_mock_test.go -package=service
type CreditLimitRepository interface {
	InsertNewCreditLimit(ctx context.Context, tx *sql.Tx, data *entity.CreditLimit) error
	FindCreditLimitByCustomerID(ctx context.Context, customerID int) (*[]entity.Limits, error)
	FindLimitByCustomerAndTenor(ctx context.Context, tx *sql.Tx, customerID, tenorMonth int) (*entity.Limits, error)
}

type CreditLimitService interface {
	GetCreditLimits(ctx context.Context, customerID int) (*[]dto.GetCreditLimitsResponse, error)
}
