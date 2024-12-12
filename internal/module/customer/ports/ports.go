package ports

import (
	"context"

	"github.com/hilmiikhsan/multifinance-service/internal/module/customer/entity"
)

type CustomerRepository interface {
	InsertNewUser(ctx context.Context, data *entity.Customer) (*entity.Customer, error)
	FindCustomerByEmail(ctx context.Context, email string) (*entity.Customer, error)
}
