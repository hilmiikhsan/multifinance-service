package ports

import (
	"context"
	"database/sql"

	"github.com/hilmiikhsan/multifinance-service/internal/module/credit_limit/entity"
)

type CreditLimitRepository interface {
	InsertNewCreditLimit(ctx context.Context, tx *sql.Tx, data *entity.CreditLimit) error
}
