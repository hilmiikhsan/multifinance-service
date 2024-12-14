package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/hilmiikhsan/multifinance-service/internal/module/credit_limit/entity"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func Test_creditLimitRepository_InsertNewCreditLimit(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mysqlDB := sqlx.NewDb(db, "mysql")

	type args struct {
		ctx   context.Context
		model *entity.CreditLimit
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
		mockFn  func(args args, mock sqlmock.Sqlmock)
	}{
		{
			name: "Insert New Credit Limit Successfully",
			args: args{
				ctx: context.Background(),
				model: &entity.CreditLimit{
					CustomerID:  1,
					TenorMonth:  12,
					LimitAmount: 50000,
				},
			},
			wantErr: false,
			mockFn: func(args args, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO credit_limits").WithArgs(
					args.model.CustomerID,
					args.model.TenorMonth,
					args.model.LimitAmount,
				).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name: "Insert New Credit Limit With Query Error",
			args: args{
				ctx: context.Background(),
				model: &entity.CreditLimit{
					CustomerID:  1,
					TenorMonth:  12,
					LimitAmount: 50000,
				},
			},
			wantErr: true,
			mockFn: func(args args, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO credit_limits").WithArgs(
					args.model.CustomerID,
					args.model.TenorMonth,
					args.model.LimitAmount,
				).WillReturnError(fmt.Errorf("insert failed"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args, mock)
			r := &creditLimitRepository{
				db: mysqlDB,
			}

			tx, err := mysqlDB.BeginTx(tt.args.ctx, nil)
			assert.NoError(t, err)

			err = r.InsertNewCreditLimit(tt.args.ctx, tx, tt.args.model)

			if tt.wantErr {
				assert.Error(t, err, "Expected error, got nil")
			} else {
				assert.NoError(t, err, "Expected no error, got: %v", err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
