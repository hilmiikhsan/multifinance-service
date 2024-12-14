package service

import (
	"context"
	"errors"
	reflect "reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/hilmiikhsan/multifinance-service/constants"
	"github.com/hilmiikhsan/multifinance-service/internal/module/auth/dto"
	creditLimitEntity "github.com/hilmiikhsan/multifinance-service/internal/module/credit_limit/entity"
	"github.com/hilmiikhsan/multifinance-service/internal/module/customer/entity"
	"github.com/hilmiikhsan/multifinance-service/pkg/jwt_handler"
	"github.com/hilmiikhsan/multifinance-service/pkg/utils"
	"github.com/jmoiron/sqlx"
	"go.uber.org/mock/gomock"
)

func Test_authService_Register(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	customerMockRepo := NewMockCustomerRepository(ctrlMock)
	creditLimitMockRepo := NewMockCreditLimitRepository(ctrlMock)

	type args struct {
		ctx context.Context
		req *dto.RegisterRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *dto.RegisterResponse
		wantErr bool
		mockFn  func(args args, dbMock sqlmock.Sqlmock)
	}{
		{
			name: "Register Success - Salary < 5M",
			args: args{
				ctx: context.Background(),
				req: &dto.RegisterRequest{
					Nik:             "123456789",
					Email:           "test@example.com",
					Password:        "testpass",
					FullName:        "Test User",
					LegalName:       "Test Legal",
					BirthPlace:      "City",
					BirthDate:       "1990-01-01",
					Salary:          4000000,
					KtpPhotoPath:    "/path/ktp.jpg",
					SelfiePhotoPath: "/path/selfie.jpg",
				},
			},
			want: &dto.RegisterResponse{
				ID:    1,
				Email: "test@example.com",
			},
			wantErr: false,
			mockFn: func(args args, dbMock sqlmock.Sqlmock) {
				dbMock.ExpectBegin()

				customerMockRepo.EXPECT().
					InsertNewUser(args.ctx, gomock.Any(), gomock.Any()).
					Return(&entity.Customer{
						ID:    1,
						Email: "test@example.com",
					}, nil)

				creditLimitMockRepo.EXPECT().
					InsertNewCreditLimit(args.ctx, gomock.Any(), &creditLimitEntity.CreditLimit{
						CustomerID: 1, TenorMonth: 1, LimitAmount: 100000,
					}).Return(nil)
				creditLimitMockRepo.EXPECT().
					InsertNewCreditLimit(args.ctx, gomock.Any(), &creditLimitEntity.CreditLimit{
						CustomerID: 1, TenorMonth: 2, LimitAmount: 200000,
					}).Return(nil)
				creditLimitMockRepo.EXPECT().
					InsertNewCreditLimit(args.ctx, gomock.Any(), &creditLimitEntity.CreditLimit{
						CustomerID: 1, TenorMonth: 3, LimitAmount: 500000,
					}).Return(nil)
				creditLimitMockRepo.EXPECT().
					InsertNewCreditLimit(args.ctx, gomock.Any(), &creditLimitEntity.CreditLimit{
						CustomerID: 1, TenorMonth: 6, LimitAmount: 700000,
					}).Return(nil)

				dbMock.ExpectCommit()
			},
		},
		{
			name: "Register Success - Salary Between 5M and 10M",
			args: args{
				ctx: context.Background(),
				req: &dto.RegisterRequest{
					Nik:             "987654321",
					Email:           "middle@example.com",
					Password:        "midpass",
					FullName:        "Middle User",
					LegalName:       "Middle Legal",
					BirthPlace:      "Town",
					BirthDate:       "1985-06-01",
					Salary:          7000000,
					KtpPhotoPath:    "/path/mid_ktp.jpg",
					SelfiePhotoPath: "/path/mid_selfie.jpg",
				},
			},
			want: &dto.RegisterResponse{
				ID:    2,
				Email: "middle@example.com",
			},
			wantErr: false,
			mockFn: func(args args, dbMock sqlmock.Sqlmock) {
				dbMock.ExpectBegin()

				customerMockRepo.EXPECT().
					InsertNewUser(args.ctx, gomock.Any(), gomock.Any()).
					Return(&entity.Customer{
						ID:    2,
						Email: "middle@example.com",
					}, nil)

				// Sesuaikan ekspektasi sesuai logika salary 5M-10M
				creditLimitMockRepo.EXPECT().
					InsertNewCreditLimit(args.ctx, gomock.Any(), &creditLimitEntity.CreditLimit{
						CustomerID: 2, TenorMonth: 1, LimitAmount: 200000,
					}).Return(nil)
				creditLimitMockRepo.EXPECT().
					InsertNewCreditLimit(args.ctx, gomock.Any(), &creditLimitEntity.CreditLimit{
						CustomerID: 2, TenorMonth: 2, LimitAmount: 400000,
					}).Return(nil)
				creditLimitMockRepo.EXPECT().
					InsertNewCreditLimit(args.ctx, gomock.Any(), &creditLimitEntity.CreditLimit{
						CustomerID: 2, TenorMonth: 3, LimitAmount: 800000,
					}).Return(nil)
				creditLimitMockRepo.EXPECT().
					InsertNewCreditLimit(args.ctx, gomock.Any(), &creditLimitEntity.CreditLimit{
						CustomerID: 2, TenorMonth: 6, LimitAmount: 1200000,
					}).Return(nil)

				dbMock.ExpectCommit()
			},
		},
		{
			name: "Error Saat Hashing Password",
			args: args{
				ctx: context.Background(),
				req: &dto.RegisterRequest{
					Password: "invalid_hash",
				},
			},
			want:    nil,
			wantErr: true,
			mockFn: func(args args, dbMock sqlmock.Sqlmock) {
				// Simulasi tidak perlu memanggil repository
			},
		},
		{
			name: "Error Saat InsertNewUser - NIK Sudah Terdaftar",
			args: args{
				ctx: context.Background(),
				req: &dto.RegisterRequest{
					Nik: "duplicate_nik",
				},
			},
			want:    nil,
			wantErr: true,
			mockFn: func(args args, dbMock sqlmock.Sqlmock) {
				dbMock.ExpectBegin()

				customerMockRepo.EXPECT().
					InsertNewUser(args.ctx, gomock.Any(), gomock.Any()).
					Return(nil, errors.New(constants.ErrNikAlreadyRegistered))

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

			s := &authService{
				db:                    mockDB,
				customerRepository:    customerMockRepo,
				creditLimitRepository: creditLimitMockRepo,
			}

			got, err := s.Register(tt.args.ctx, tt.args.req)

			if (err != nil) != tt.wantErr {
				t.Errorf("authService.Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("authService.Register() = %v, want %v", got, tt.want)
			}

			if err := dbMock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %v", err)
			}
		})
	}
}

func Test_authService_Login(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	customerMockRepo := NewMockCustomerRepository(ctrlMock)
	mockJWT := NewMockJWT(ctrlMock)

	password, _ := utils.HashPassword("password123")

	type args struct {
		ctx context.Context
		req *dto.LoginRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *dto.LoginResponse
		wantErr bool
		mockFn  func(args args)
	}{
		{
			name: "Login Success",
			args: args{
				ctx: context.Background(),
				req: &dto.LoginRequest{
					Email:    "test@example.com",
					Password: "password123",
				},
			},
			want: &dto.LoginResponse{
				ID:           1,
				Email:        "test@example.com",
				FullName:     "Test User",
				Token:        "access-token",
				RefreshToken: "refresh-token",
			},
			wantErr: false,
			mockFn: func(args args) {
				customerMockRepo.EXPECT().
					FindCustomerByEmail(args.ctx, args.req.Email).
					Return(&entity.Customer{
						ID:       1,
						Nik:      "123456789",
						Email:    "test@example.com",
						FullName: "Test User",
						Password: password,
					}, nil)

				mockJWT.EXPECT().
					GenerateTokenString(args.ctx, jwt_handler.CostumClaimsPayload{
						CustomerID: 1,
						Nik:        "123456789",
						Email:      "test@example.com",
						FullName:   "Test User",
						TokenType:  constants.AccessTokenType,
					}).
					Return("access-token", nil)

				mockJWT.EXPECT().
					GenerateTokenString(args.ctx, jwt_handler.CostumClaimsPayload{
						CustomerID: 1,
						Nik:        "123456789",
						Email:      "test@example.com",
						FullName:   "Test User",
						TokenType:  constants.RefreshTokenType,
					}).
					Return("refresh-token", nil)
			},
		},
		{
			name: "Email Not Found",
			args: args{
				ctx: context.Background(),
				req: &dto.LoginRequest{
					Email:    "notfound@example.com",
					Password: "password123",
				},
			},
			want:    nil,
			wantErr: true,
			mockFn: func(args args) {
				customerMockRepo.EXPECT().
					FindCustomerByEmail(args.ctx, args.req.Email).
					Return(nil, nil)
			},
		},
		{
			name: "Incorrect Password",
			args: args{
				ctx: context.Background(),
				req: &dto.LoginRequest{
					Email:    "test@example.com",
					Password: "wrongpassword",
				},
			},
			want:    nil,
			wantErr: true,
			mockFn: func(args args) {
				customerMockRepo.EXPECT().
					FindCustomerByEmail(args.ctx, args.req.Email).
					Return(&entity.Customer{
						ID:       1,
						Nik:      "123456789",
						Email:    "test@example.com",
						FullName: "Test User",
						Password: password,
					}, nil)
			},
		},
		{
			name: "Error Generating Access Token",
			args: args{
				ctx: context.Background(),
				req: &dto.LoginRequest{
					Email:    "test@example.com",
					Password: "password123",
				},
			},
			want:    nil,
			wantErr: true,
			mockFn: func(args args) {
				customerMockRepo.EXPECT().
					FindCustomerByEmail(args.ctx, args.req.Email).
					Return(&entity.Customer{
						ID:       1,
						Nik:      "123456789",
						Email:    "test@example.com",
						FullName: "Test User",
						Password: password,
					}, nil)

				mockJWT.EXPECT().
					GenerateTokenString(args.ctx, jwt_handler.CostumClaimsPayload{
						CustomerID: 1,
						Nik:        "123456789",
						Email:      "test@example.com",
						FullName:   "Test User",
						TokenType:  constants.AccessTokenType,
					}).
					Return("", errors.New("failed to generate token"))
			},
		},
		{
			name: "Error Generating Refresh Token",
			args: args{
				ctx: context.Background(),
				req: &dto.LoginRequest{
					Email:    "test@example.com",
					Password: "password123",
				},
			},
			want:    nil,
			wantErr: true,
			mockFn: func(args args) {
				customerMockRepo.EXPECT().
					FindCustomerByEmail(args.ctx, args.req.Email).
					Return(&entity.Customer{
						ID:       1,
						Nik:      "123456789",
						Email:    "test@example.com",
						FullName: "Test User",
						Password: password,
					}, nil)

				mockJWT.EXPECT().
					GenerateTokenString(args.ctx, jwt_handler.CostumClaimsPayload{
						CustomerID: 1,
						Nik:        "123456789",
						Email:      "test@example.com",
						FullName:   "Test User",
						TokenType:  constants.AccessTokenType,
					}).
					Return("access-token", nil)

				mockJWT.EXPECT().
					GenerateTokenString(args.ctx, jwt_handler.CostumClaimsPayload{
						CustomerID: 1,
						Nik:        "123456789",
						Email:      "test@example.com",
						FullName:   "Test User",
						TokenType:  constants.RefreshTokenType,
					}).
					Return("", errors.New("failed to generate token"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)

			s := &authService{
				customerRepository: customerMockRepo,
				jwt:                mockJWT,
			}

			got, err := s.Login(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("authService.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("authService.Login() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_authService_RefreshToken(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockJWT := NewMockJWT(ctrlMock)

	type args struct {
		ctx         context.Context
		accessToken string
	}
	tests := []struct {
		name    string
		args    args
		want    *dto.RefreshTokenResponse
		wantErr bool
		mockFn  func(args args)
	}{
		{
			name: "RefreshToken Success",
			args: args{
				ctx:         context.Background(),
				accessToken: "valid-access-token",
			},
			want: &dto.RefreshTokenResponse{
				Token: "new-access-token",
			},
			wantErr: false,
			mockFn: func(args args) {
				mockJWT.EXPECT().
					ParseTokenString(args.ctx, args.accessToken).
					Return(&jwt_handler.CustomClaims{
						Nik:      "123456789",
						Email:    "test@example.com",
						FullName: "Test User",
					}, nil)

				mockJWT.EXPECT().
					GenerateTokenString(args.ctx, jwt_handler.CostumClaimsPayload{
						Nik:       "123456789",
						Email:     "test@example.com",
						FullName:  "Test User",
						TokenType: constants.AccessTokenType,
					}).
					Return("new-access-token", nil)
			},
		},
		{
			name: "Invalid Access Token",
			args: args{
				ctx:         context.Background(),
				accessToken: "invalid-access-token",
			},
			want:    nil,
			wantErr: true,
			mockFn: func(args args) {
				mockJWT.EXPECT().
					ParseTokenString(args.ctx, args.accessToken).
					Return(nil, errors.New("invalid token"))
			},
		},
		{
			name: "Error Generating New Token",
			args: args{
				ctx:         context.Background(),
				accessToken: "valid-access-token",
			},
			want:    nil,
			wantErr: true,
			mockFn: func(args args) {
				mockJWT.EXPECT().
					ParseTokenString(args.ctx, args.accessToken).
					Return(&jwt_handler.CustomClaims{
						Nik:      "123456789",
						Email:    "test@example.com",
						FullName: "Test User",
					}, nil)

				mockJWT.EXPECT().
					GenerateTokenString(args.ctx, jwt_handler.CostumClaimsPayload{
						Nik:       "123456789",
						Email:     "test@example.com",
						FullName:  "Test User",
						TokenType: constants.AccessTokenType,
					}).
					Return("", errors.New("failed to generate token"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)

			s := &authService{
				jwt: mockJWT,
			}

			got, err := s.RefreshToken(tt.args.ctx, tt.args.accessToken)

			if (err != nil) != tt.wantErr {
				t.Errorf("authService.RefreshToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("authService.RefreshToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
