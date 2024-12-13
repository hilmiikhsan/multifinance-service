package ports

import (
	"context"
	"database/sql"

	"github.com/hilmiikhsan/multifinance-service/internal/module/customer/dto"
	"github.com/hilmiikhsan/multifinance-service/internal/module/customer/entity"
)

type CustomerRepository interface {
	InsertNewUser(ctx context.Context, tx *sql.Tx, data *entity.Customer) (*entity.Customer, error)
	FindCustomerByEmail(ctx context.Context, email string) (*entity.Customer, error)
	FindCustomerByID(ctx context.Context, id int) (*entity.Customer, error)
}

type CustomerService interface {
	GetCustomerProfile(ctx context.Context, id int) (*dto.GetCustomerProfileResponse, error)
}
