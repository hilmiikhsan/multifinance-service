package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/hilmiikhsan/multifinance-service/internal/module/transaction/entity"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func Test_transactionRepository_InsertNewTransaction(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mysqlDB := sqlx.NewDb(db, "mysql")

	type args struct {
		ctx   context.Context
		model *entity.Transaction
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		mockFn  func(args args, mock sqlmock.Sqlmock)
	}{
		{
			name: "Insert New Transaction Successfully",
			args: args{
				ctx: context.Background(),
				model: &entity.Transaction{
					CustomerID:        1,
					ContractNumber:    "123456",
					OnTheRoadPrice:    500000,
					AdminFee:          5000,
					InstallmentAmount: 50000,
					InterestAmount:    5000,
					AssetName:         "Yamaha NMAX",
				},
			},
			wantErr: false,
			mockFn: func(args args, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO transactions").WithArgs(
					args.model.CustomerID,
					args.model.ContractNumber,
					args.model.OnTheRoadPrice,
					args.model.AdminFee,
					args.model.InstallmentAmount,
					args.model.InterestAmount,
					args.model.AssetName,
				).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name: "Insert New Transaction With Query Error",
			args: args{
				ctx: context.Background(),
				model: &entity.Transaction{
					CustomerID:        1,
					ContractNumber:    "123456",
					OnTheRoadPrice:    500000,
					AdminFee:          5000,
					InstallmentAmount: 50000,
					InterestAmount:    5000,
					AssetName:         "Yamaha NMAX",
				},
			},
			wantErr: true,
			mockFn: func(args args, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO transactions").WithArgs(
					args.model.CustomerID,
					args.model.ContractNumber,
					args.model.OnTheRoadPrice,
					args.model.AdminFee,
					args.model.InstallmentAmount,
					args.model.InterestAmount,
					args.model.AssetName,
				).WillReturnError(fmt.Errorf("insert failed"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args, mock)
			r := &transactionRepository{
				db: mysqlDB,
			}

			tx, err := mysqlDB.BeginTx(tt.args.ctx, nil)
			assert.NoError(t, err)

			err = r.InsertNewTransaction(tt.args.ctx, tx, tt.args.model)

			if tt.wantErr {
				assert.Error(t, err, "Expected error, got nil")
			} else {
				assert.NoError(t, err, "Expected no error, got: %v", err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func Test_transactionRepository_FindTransactionByIdAndCustomerID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mysqlDB := sqlx.NewDb(db, "mysql")

	type args struct {
		ctx        context.Context
		id         int
		customerID int
	}
	tests := []struct {
		name    string
		args    args
		want    *entity.Transaction
		wantErr bool
		mockFn  func(args args, mock sqlmock.Sqlmock)
	}{
		{
			name: "Find Transaction By ID and Customer ID Successfully",
			args: args{
				ctx:        context.Background(),
				id:         1,
				customerID: 1,
			},
			want: &entity.Transaction{
				ID:                1,
				CustomerID:        1,
				ContractNumber:    "123456",
				OnTheRoadPrice:    500000,
				AdminFee:          5000,
				InstallmentAmount: 50000,
				InterestAmount:    5000,
				AssetName:         "Yamaha NMAX",
			},
			wantErr: false,
			mockFn: func(args args, mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "customer_id", "contract_number", "on_the_road_price", "admin_fee", "installment_amount", "interest_amount", "asset_name"}).
					AddRow(1, 1, "123456", 500000, 5000, 50000, 5000, "Yamaha NMAX")

				mock.ExpectQuery("SELECT (.+) FROM transactions").WithArgs(args.id, args.customerID).WillReturnRows(rows)
			},
		},
		{
			name: "Find Transaction By ID and Customer ID With Query Error",
			args: args{
				ctx:        context.Background(),
				id:         1,
				customerID: 1,
			},
			want:    nil,
			wantErr: true,
			mockFn: func(args args, mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM transactions").WithArgs(args.id, args.customerID).WillReturnError(fmt.Errorf("query failed"))
			},
		},
		{
			name: "Find Transaction By ID and Customer ID With No Rows",
			args: args{
				ctx:        context.Background(),
				id:         1,
				customerID: 1,
			},
			want:    nil,
			wantErr: true,
			mockFn: func(args args, mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "customer_id", "contract_number", "on_the_road_price", "admin_fee", "installment_amount", "interest_amount", "asset_name"})

				mock.ExpectQuery("SELECT (.+) FROM transactions").WithArgs(args.id, args.customerID).WillReturnRows(rows)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args, mock)
			r := &transactionRepository{
				db: mysqlDB,
			}
			got, err := r.FindTransactionByIdAndCustomerID(tt.args.ctx, tt.args.id, tt.args.customerID)

			assert.Equal(t, tt.wantErr, err != nil, "error state mismatch")
			assert.Equal(t, tt.want, got, "result mismatch")
		})
	}

	assert.NoError(t, mock.ExpectationsWereMet())
}
