package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/hilmiikhsan/multifinance-service/constants"
	creditLimitEntity "github.com/hilmiikhsan/multifinance-service/internal/module/credit_limit/entity"
	"github.com/hilmiikhsan/multifinance-service/internal/module/transaction/dto"
	"github.com/hilmiikhsan/multifinance-service/internal/module/transaction/entity"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func Test_transactionService_CreateTransaction(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockTransactionRepo := NewMockTransactionRepository(ctrlMock)
	mockCreditLimitRepo := NewMockCreditLimitRepository(ctrlMock)

	type args struct {
		ctx context.Context
		req *dto.CreateTransactionRequest
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
		mockFn  func(args args, dbMock sqlmock.Sqlmock)
	}{
		{
			name: "CreateTransaction Success",
			args: args{
				ctx: context.Background(),
				req: &dto.CreateTransactionRequest{
					CustomerID:        1,
					OnTheRoadPrice:    500000,
					InstallmentAmount: 50000,
					InterestAmount:    5000,
					AssetName:         "Yamaha NMAX",
					TenorMonth:        12,
				},
			},
			wantErr: false,
			mockFn: func(args args, dbMock sqlmock.Sqlmock) {
				dbMock.ExpectBegin()

				mockCreditLimitRepo.EXPECT().FindLimitByCustomerAndTenor(args.ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(&creditLimitEntity.Limits{
					LimitAmount: 1000000,
				}, nil)

				mockTransactionRepo.EXPECT().InsertNewTransaction(args.ctx, gomock.Any(), gomock.Any()).Return(nil)

				dbMock.ExpectCommit()
			},
		},
		{
			name: "CreateTransaction Failed - Limit Not Found",
			args: args{
				ctx: context.Background(),
				req: &dto.CreateTransactionRequest{
					CustomerID:        1,
					OnTheRoadPrice:    500000,
					InstallmentAmount: 50000,
					InterestAmount:    5000,
					AssetName:         "Yamaha NMAX",
					TenorMonth:        12,
				},
			},
			wantErr: true,
			mockFn: func(args args, dbMock sqlmock.Sqlmock) {
				dbMock.ExpectBegin()

				mockCreditLimitRepo.EXPECT().
					FindLimitByCustomerAndTenor(args.ctx, gomock.Any(), args.req.CustomerID, args.req.TenorMonth).
					Return(nil, errors.New(constants.ErrInvalidOrCreditLimit))

				dbMock.ExpectRollback()
			},
		},
		{
			name: "CreateTransaction Failed - Insert Transaction Error",
			args: args{
				ctx: context.Background(),
				req: &dto.CreateTransactionRequest{
					CustomerID:        1,
					OnTheRoadPrice:    500000,
					InstallmentAmount: 50000,
					InterestAmount:    5000,
					AssetName:         "Yamaha NMAX",
					TenorMonth:        12,
				},
			},
			wantErr: true,
			mockFn: func(args args, dbMock sqlmock.Sqlmock) {
				dbMock.ExpectBegin()

				mockCreditLimitRepo.EXPECT().FindLimitByCustomerAndTenor(args.ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(&creditLimitEntity.Limits{
					LimitAmount: 1000000,
				}, nil)

				mockTransactionRepo.EXPECT().InsertNewTransaction(args.ctx, gomock.Any(), gomock.Any()).Return(errors.New(constants.ErrInternalServerError))

				dbMock.ExpectCommit()
			},
		},
		{
			name: "CreateTransaction Failed - Begin Transaction Error",
			args: args{
				ctx: context.Background(),
				req: &dto.CreateTransactionRequest{
					CustomerID:        1,
					OnTheRoadPrice:    500000,
					InstallmentAmount: 50000,
					InterestAmount:    5000,
					AssetName:         "Yamaha NMAX",
					TenorMonth:        12,
				},
			},
			wantErr: true,
			mockFn: func(args args, dbMock sqlmock.Sqlmock) {
				dbMock.ExpectBegin().WillReturnError(errors.New(constants.ErrInternalServerError))
			},
		},
		{
			name: "CreateTransaction Failed - Commit Transaction Error",
			args: args{
				ctx: context.Background(),
				req: &dto.CreateTransactionRequest{
					CustomerID:        1,
					OnTheRoadPrice:    500000,
					InstallmentAmount: 50000,
					InterestAmount:    5000,
					AssetName:         "Yamaha NMAX",
					TenorMonth:        12,
				},
			},
			wantErr: true,
			mockFn: func(args args, dbMock sqlmock.Sqlmock) {
				dbMock.ExpectBegin()

				mockCreditLimitRepo.EXPECT().FindLimitByCustomerAndTenor(args.ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(&creditLimitEntity.Limits{
					LimitAmount: 1000000,
				}, nil)

				mockTransactionRepo.EXPECT().InsertNewTransaction(args.ctx, gomock.Any(), gomock.Any()).Return(nil)

				dbMock.ExpectCommit().WillReturnError(errors.New(constants.ErrInternalServerError))
			},
		},
		{
			name: "CreateTransaction Failed - Rollback Transaction Error",
			args: args{
				ctx: context.Background(),
				req: &dto.CreateTransactionRequest{
					CustomerID:        1,
					OnTheRoadPrice:    500000,
					InstallmentAmount: 50000,
					InterestAmount:    5000,
					AssetName:         "Yamaha NMAX",
					TenorMonth:        12,
				},
			},
			wantErr: true,
			mockFn: func(args args, dbMock sqlmock.Sqlmock) {
				dbMock.ExpectBegin()

				mockCreditLimitRepo.EXPECT().
					FindLimitByCustomerAndTenor(args.ctx, gomock.Any(), args.req.CustomerID, args.req.TenorMonth).
					Return(nil, errors.New(constants.ErrInvalidOrCreditLimit))

				dbMock.ExpectRollback().WillReturnError(errors.New(constants.ErrInternalServerError))
			},
		},
		{
			name: "CreateTransaction Failed - Find Limit By Customer And Tenor Error",
			args: args{
				ctx: context.Background(),
				req: &dto.CreateTransactionRequest{
					CustomerID:        1,
					OnTheRoadPrice:    50000000,
					InstallmentAmount: 50000,
					InterestAmount:    5000,
					AssetName:         "Yamaha NMAX",
					TenorMonth:        12,
				},
			},
			wantErr: true,
			mockFn: func(args args, dbMock sqlmock.Sqlmock) {
				dbMock.ExpectBegin()

				mockCreditLimitRepo.EXPECT().FindLimitByCustomerAndTenor(args.ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(&creditLimitEntity.Limits{
					LimitAmount: 1000000,
				}, nil)

				if args.req.OnTheRoadPrice > 1000000 {
					assert.Error(t, errors.New(constants.ErrOnTheRoadPriceExceedLimit), "expected an error but got none")
				}

				dbMock.ExpectRollback()
			},
		},
		{
			name: "CreateTransaction Failed - On the road price exceeds credit limit",
			args: args{
				ctx: context.Background(),
				req: &dto.CreateTransactionRequest{
					CustomerID:        1,
					OnTheRoadPrice:    500000,
					InstallmentAmount: 50000,
					InterestAmount:    5000,
					AssetName:         "Yamaha NMAX",
					TenorMonth:        12,
				},
			},
			wantErr: true,
			mockFn: func(args args, dbMock sqlmock.Sqlmock) {
				dbMock.ExpectBegin()

				mockCreditLimitRepo.EXPECT().
					FindLimitByCustomerAndTenor(args.ctx, gomock.Any(), args.req.CustomerID, args.req.TenorMonth).
					Return(&creditLimitEntity.Limits{
						LimitAmount: 100000,
					}, nil)

				dbMock.ExpectRollback()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, dbMock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create sqlmock: %v", err)
			}
			defer db.Close()

			mockDB := sqlx.NewDb(db, "mysql")

			tt.mockFn(tt.args, dbMock)

			s := &transactionService{
				db:                    mockDB,
				transactionRepository: mockTransactionRepo,
				creditLimitRepository: mockCreditLimitRepo,
			}
			err = s.CreateTransaction(tt.args.ctx, tt.args.req)

			if tt.wantErr {
				assert.Error(t, err, "expected an error but got none")
			} else {
				assert.NoError(t, err, "did not expect an error but got one")
			}
		})
	}
}

func Test_transactionService_GetDetailTransaction(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockRepo := NewMockTransactionRepository(ctrlMock)

	type args struct {
		ctx        context.Context
		id         int
		customerID int
	}
	tests := []struct {
		name    string
		args    args
		want    *dto.GetDetailTransactionResponse
		wantErr bool
		mockFn  func(args args)
	}{
		{
			name: "GetDetailTransaction Success",
			args: args{
				ctx:        context.Background(),
				id:         1,
				customerID: 1,
			},
			want: &dto.GetDetailTransactionResponse{
				ID:                1,
				CustomerID:        1,
				OnTheRoadPrice:    500000,
				InstallmentAmount: 50000,
				InterestAmount:    5000,
				AssetName:         "Yamaha NMAX",
				CreatedAt:         time.Now().Format("2006-01-02 15:04:05"),
			},
			wantErr: false,
			mockFn: func(args args) {
				mockRepo.EXPECT().FindTransactionByIdAndCustomerID(gomock.Any(), args.id, args.customerID).Return(&entity.Transaction{
					ID:                1,
					CustomerID:        1,
					OnTheRoadPrice:    500000,
					InstallmentAmount: 50000,
					InterestAmount:    5000,
					AssetName:         "Yamaha NMAX",
					CreatedAt:         time.Now(),
				}, nil)
			},
		},
		{
			name: "GetDetailTransaction Failed - Transaction Not Found",
			args: args{
				ctx:        context.Background(),
				id:         1,
				customerID: 1,
			},
			want:    nil,
			wantErr: true,
			mockFn: func(args args) {
				mockRepo.EXPECT().FindTransactionByIdAndCustomerID(gomock.Any(), args.id, args.customerID).Return(nil, errors.New(constants.ErrTransactionNotFound))
			},
		},
		{
			name: "GetDetailTransaction Failed - Internal Server Error",
			args: args{
				ctx:        context.Background(),
				id:         1,
				customerID: 1,
			},
			want:    nil,
			wantErr: true,
			mockFn: func(args args) {
				mockRepo.EXPECT().FindTransactionByIdAndCustomerID(gomock.Any(), args.id, args.customerID).Return(nil, errors.New(constants.ErrInternalServerError))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)

			s := &transactionService{
				transactionRepository: mockRepo,
			}

			got, err := s.GetDetailTransaction(tt.args.ctx, tt.args.id, tt.args.customerID)

			if tt.wantErr {
				assert.Error(t, err, "expected an error but got none")
			} else {
				assert.NoError(t, err, "did not expect an error but got one")
			}

			assert.Equal(t, tt.want, got, "unexpected result from GetCustomerProfile")
		})
	}
}
