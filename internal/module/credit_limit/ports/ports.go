package ports

import (
	"context"
	"database/sql"

	"github.com/hilmiikhsan/multifinance-service/internal/module/credit_limit/dto"
	"github.com/hilmiikhsan/multifinance-service/internal/module/credit_limit/entity"
)

type CreditLimitRepository interface {
	InsertNewCreditLimit(ctx context.Context, tx *sql.Tx, data *entity.CreditLimit) error
	FindCreditLimitByCustomerID(ctx context.Context, customerID int) (*[]entity.Limits, error)
}

type CreditLimitService interface {
	GetCreditLimits(ctx context.Context, customerID int) (*[]dto.GetCreditLimitsResponse, error)
}
