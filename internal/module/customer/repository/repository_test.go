package repository

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	creditLimitEntity "github.com/hilmiikhsan/multifinance-service/internal/module/credit_limit/entity"
	"github.com/hilmiikhsan/multifinance-service/internal/module/customer/entity"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func Test_customerRepository_InsertNewUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mysqlDB := sqlx.NewDb(db, "mysql")

	birthDate, err := time.Parse("2006-01-02", "1990-01-01")
	assert.NoError(t, err)

	type args struct {
		ctx   context.Context
		model *entity.Customer
	}
	tests := []struct {
		name    string
		args    args
		want    *entity.Customer
		wantErr bool
		mockFn  func(args args, mock sqlmock.Sqlmock)
	}{
		{
			name: "Insert New User Successfully",
			args: args{
				ctx: context.Background(),
				model: &entity.Customer{
					Nik:             "123456789",
					Email:           "test@domain.com",
					Password:        "hashed_password",
					FullName:        "Test User",
					LegalName:       "Test User",
					BirthPlace:      "City",
					BirthDate:       birthDate,
					Salary:          10000,
					KtpPhotoPath:    "/path/to/ktp/photo",
					SelfiePhotoPath: "/path/to/selfie/photo",
				},
			},
			wantErr: false,
			mockFn: func(args args, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO customers").WithArgs(
					args.model.Nik,
					args.model.Email,
					args.model.Password,
					args.model.FullName,
					args.model.LegalName,
					args.model.BirthPlace,
					args.model.BirthDate,
					args.model.Salary,
					args.model.KtpPhotoPath,
					args.model.SelfiePhotoPath,
				).WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectQuery("SELECT id, email FROM customers WHERE id = ?").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "email"}).AddRow(1, args.model.Email))
				mock.ExpectCommit()
			},
		},
		{
			name: "Insert New User - Unique Constraint Violation (NIK)",
			args: args{
				ctx: context.Background(),
				model: &entity.Customer{
					Nik:   "123456789",
					Email: "new@domain.com",
				},
			},
			wantErr: true,
			mockFn: func(args args, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO customers").WithArgs(
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
				).WillReturnError(fmt.Errorf("Error 1062: Duplicate entry '123456789' for key 'nik'"))
				mock.ExpectRollback()
			},
		},
		{
			name: "Insert New User - Unique Constraint Violation (Email)",
			args: args{
				ctx: context.Background(),
				model: &entity.Customer{
					Nik:   "987654321",
					Email: "existing@domain.com",
				},
			},
			wantErr: true,
			mockFn: func(args args, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO customers").WithArgs(
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
				).WillReturnError(fmt.Errorf("Error 1062: Duplicate entry 'existing@domain.com' for key 'email'"))
				mock.ExpectRollback()
			},
		},
		{
			name: "Insert New User - Error Getting Last Insert ID",
			args: args{
				ctx: context.Background(),
				model: &entity.Customer{
					Nik:   "123456789",
					Email: "test@domain.com",
				},
			},
			wantErr: true,
			mockFn: func(args args, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO customers").WithArgs(
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
				).WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("Error getting last insert ID")))
				mock.ExpectRollback()
			},
		},
		{
			name: "Insert New User - Error Querying User Details",
			args: args{
				ctx: context.Background(),
				model: &entity.Customer{
					Nik:   "123456789",
					Email: "test@domain.com",
				},
			},
			wantErr: true,
			mockFn: func(args args, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO customers").WithArgs(
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
				).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectQuery("SELECT id, email FROM customers WHERE id = ?").
					WithArgs(1).
					WillReturnError(fmt.Errorf("Error querying user details"))
				mock.ExpectRollback()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args, mock)
			r := &customerRepository{
				db: mysqlDB,
			}

			tx, err := mysqlDB.BeginTx(tt.args.ctx, nil)
			assert.NoError(t, err)

			got, err := r.InsertNewUser(tt.args.ctx, tx, tt.args.model)

			if tt.wantErr {
				assert.Error(t, err)
				switch tt.name {
				case "Insert New User - Unique Constraint Violation (NIK)":
					assert.Contains(t, err.Error(), "Duplicate entry")
					assert.Contains(t, err.Error(), "for key 'nik'")
				case "Insert New User - Unique Constraint Violation (Email)":
					assert.Contains(t, err.Error(), "Duplicate entry")
					assert.Contains(t, err.Error(), "for key 'email'")
				case "Insert New User - Error Getting Last Insert ID":
					assert.Contains(t, err.Error(), "Internal server error")
				case "Insert New User - Error Querying User Details":
					assert.Contains(t, err.Error(), "Internal server error")
				}
				_ = tx.Rollback()
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				if got != nil {
					assert.Equal(t, int64(1), got.ID)
					assert.Equal(t, tt.args.model.Email, got.Email)
				}
				err = tx.Commit()
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func Test_customerRepository_FindCustomerByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mysqlDB := sqlx.NewDb(db, "mysql")

	type args struct {
		ctx   context.Context
		email string
	}
	tests := []struct {
		name    string
		args    args
		want    *entity.Customer
		wantErr bool
		mockFn  func(args args, mock sqlmock.Sqlmock)
	}{
		{
			name: "Customer found successfully",
			args: args{
				ctx:   context.Background(),
				email: "test@example.com",
			},
			want: &entity.Customer{
				ID:    1,
				Email: "test@example.com",
			},
			wantErr: false,
			mockFn: func(args args, mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "email"}).
					AddRow(1, "test@example.com")
				mock.ExpectQuery(regexp.QuoteMeta(queryFindCustomerByEmail)).
					WithArgs(args.email).
					WillReturnRows(rows)
			},
		},
		{
			name: "Customer not found",
			args: args{
				ctx:   context.Background(),
				email: "nonexistent@example.com",
			},
			want:    nil,
			wantErr: false,
			mockFn: func(args args, mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(queryFindCustomerByEmail)).
					WithArgs(args.email).
					WillReturnError(sql.ErrNoRows)
			},
		},
		{
			name: "Database error",
			args: args{
				ctx:   context.Background(),
				email: "error@example.com",
			},
			want:    nil,
			wantErr: true,
			mockFn: func(args args, mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(queryFindCustomerByEmail)).
					WithArgs(args.email).
					WillReturnError(fmt.Errorf("database error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args, mock)
			r := &customerRepository{
				db: mysqlDB,
			}
			got, err := r.FindCustomerByEmail(tt.args.ctx, tt.args.email)

			assert.Equal(t, tt.wantErr, err != nil, "error state mismatch")
			assert.Equal(t, tt.want, got, "result mismatch")
		})
	}

	assert.NoError(t, mock.ExpectationsWereMet())
}

func compareCustomersIgnoringTimestamps(t *testing.T, got, want *entity.Customer) {
	t.Helper()

	// Compare non-timestamp fields
	assert.Equal(t, want.ID, got.ID)
	assert.Equal(t, want.Nik, got.Nik)
	assert.Equal(t, want.Email, got.Email)
	assert.Equal(t, want.FullName, got.FullName)
	assert.Equal(t, want.LegalName, got.LegalName)
	assert.Equal(t, want.BirthPlace, got.BirthPlace)
	assert.Equal(t, want.BirthDate, got.BirthDate)
	assert.Equal(t, want.Salary, got.Salary)
	assert.Equal(t, want.KtpPhotoPath, got.KtpPhotoPath)
	assert.Equal(t, want.SelfiePhotoPath, got.SelfiePhotoPath)

	// Compare Limits
	assert.Equal(t, len(want.Limits), len(got.Limits))
	for i := range want.Limits {
		assert.Equal(t, want.Limits[i], got.Limits[i])
	}

	// Check that timestamps are not zero, but don't compare their exact values
	assert.False(t, got.CreatedAt.IsZero())
	assert.False(t, got.UpdatedAt.IsZero())
}

func Test_customerRepository_FindCustomerByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mysqlDB := sqlx.NewDb(db, "mysql")

	type args struct {
		ctx context.Context
		id  int
	}
	tests := []struct {
		name    string
		args    args
		want    *entity.Customer
		wantErr bool
		mockFn  func(args args, mock sqlmock.Sqlmock)
	}{
		{
			name: "Customer found with limits",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			want: &entity.Customer{
				ID:              1,
				Nik:             "1234567890",
				Email:           "test@example.com",
				FullName:        "Test User",
				LegalName:       "Test User Legal",
				BirthPlace:      "Test City",
				BirthDate:       time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
				Salary:          5000.00,
				KtpPhotoPath:    "path/to/ktp.jpg",
				SelfiePhotoPath: "path/to/selfie.jpg",
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
				Limits: []creditLimitEntity.Limits{
					{TenorMonth: 6, LimitAmount: 1000.00},
					{TenorMonth: 12, LimitAmount: 2000.00},
				},
			},
			wantErr: false,
			mockFn: func(args args, mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "nik", "email", "full_name", "legal_name", "birth_place", "birth_date",
					"salary", "ktp_photo_path", "selfie_photo_path", "created_at", "updated_at",
					"tenor_month", "limit_amount",
				}).AddRow(
					1, "1234567890", "test@example.com", "Test User", "Test User Legal", "Test City",
					time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC), 5000.00, "path/to/ktp.jpg", "path/to/selfie.jpg",
					time.Now(), time.Now(), 6, 1000.00,
				).AddRow(
					1, "1234567890", "test@example.com", "Test User", "Test User Legal", "Test City",
					time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC), 5000.00, "path/to/ktp.jpg", "path/to/selfie.jpg",
					time.Now(), time.Now(), 12, 2000.00,
				)

				mock.ExpectQuery(regexp.QuoteMeta(queryFindCustomerByID)).
					WithArgs(args.id).
					WillReturnRows(rows)
			},
		},
		{
			name: "Customer found without limits",
			args: args{
				ctx: context.Background(),
				id:  2,
			},
			want: &entity.Customer{
				ID:              2,
				Nik:             "0987654321",
				Email:           "test2@example.com",
				FullName:        "Test User 2",
				LegalName:       "Test User 2 Legal",
				BirthPlace:      "Test City 2",
				BirthDate:       time.Date(1995, 1, 1, 0, 0, 0, 0, time.UTC),
				Salary:          6000.00,
				KtpPhotoPath:    "path/to/ktp2.jpg",
				SelfiePhotoPath: "path/to/selfie2.jpg",
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
				Limits:          []creditLimitEntity.Limits{},
			},
			wantErr: false,
			mockFn: func(args args, mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "nik", "email", "full_name", "legal_name", "birth_place", "birth_date",
					"salary", "ktp_photo_path", "selfie_photo_path", "created_at", "updated_at",
					"tenor_month", "limit_amount",
				}).AddRow(
					2, "0987654321", "test2@example.com", "Test User 2", "Test User 2 Legal", "Test City 2",
					time.Date(1995, 1, 1, 0, 0, 0, 0, time.UTC), 6000.00, "path/to/ktp2.jpg", "path/to/selfie2.jpg",
					time.Now(), time.Now(), nil, nil,
				)

				mock.ExpectQuery(regexp.QuoteMeta(queryFindCustomerByID)).
					WithArgs(args.id).
					WillReturnRows(rows)
			},
		},
		{
			name: "Customer not found",
			args: args{
				ctx: context.Background(),
				id:  999,
			},
			want:    nil,
			wantErr: true,
			mockFn: func(args args, mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(queryFindCustomerByID)).
					WithArgs(args.id).
					WillReturnError(sql.ErrNoRows)
			},
		},
		{
			name: "Database error",
			args: args{
				ctx: context.Background(),
				id:  3,
			},
			want:    nil,
			wantErr: true,
			mockFn: func(args args, mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(queryFindCustomerByID)).
					WithArgs(args.id).
					WillReturnError(fmt.Errorf("database error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args, mock)
			r := &customerRepository{
				db: mysqlDB,
			}
			got, err := r.FindCustomerByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("customerRepository.FindCustomerByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil {
				compareCustomersIgnoringTimestamps(t, got, tt.want)
			} else {
				assert.Nil(t, got)
			}
		})
	}
}
