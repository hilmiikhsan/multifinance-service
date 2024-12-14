package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/hilmiikhsan/multifinance-service/constants"
	"github.com/hilmiikhsan/multifinance-service/internal/module/credit_limit/dto"
	creditLimitEntity "github.com/hilmiikhsan/multifinance-service/internal/module/credit_limit/entity"
	customerDto "github.com/hilmiikhsan/multifinance-service/internal/module/customer/dto"
	"github.com/hilmiikhsan/multifinance-service/internal/module/customer/entity"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_customerService_GetCustomerProfile(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockRepo := NewMockCustomerRepository(ctrlMock)

	type args struct {
		ctx context.Context
		id  int
	}
	tests := []struct {
		name    string
		args    args
		want    *customerDto.GetCustomerProfileResponse
		wantErr bool
		mockFn  func(args args)
	}{
		{
			name: "GetCustomerProfile Success",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			want: &customerDto.GetCustomerProfileResponse{
				ID:              1,
				Nik:             "123456789",
				FullName:        "Test User",
				LegalName:       "Test User",
				BirthPlace:      "City",
				BirthDate:       "1990-01-01",
				Salary:          10000,
				KtpPhotoPath:    "/path/to/ktp/photo",
				SelfiePhotoPath: "/path/to/selfie/photo",
				Limits: []dto.CreditLimit{
					{Tenor: 12, LimitAmount: 50000},
					{Tenor: 24, LimitAmount: 100000},
				},
				CreatedAt: "2024-01-01 10:00:00",
				UpdatedAt: "2024-01-01 10:00:00",
			},
			wantErr: false,
			mockFn: func(args args) {
				mockRepo.EXPECT().FindCustomerByID(gomock.Any(), args.id).Return(&entity.Customer{
					ID:              1,
					Nik:             "123456789",
					FullName:        "Test User",
					LegalName:       "Test User",
					BirthPlace:      "City",
					BirthDate:       parseDate("1990-01-01"),
					Salary:          10000,
					KtpPhotoPath:    "/path/to/ktp/photo",
					SelfiePhotoPath: "/path/to/selfie/photo",
					Limits: []creditLimitEntity.Limits{
						{TenorMonth: 12, LimitAmount: 50000},
						{TenorMonth: 24, LimitAmount: 100000},
					},
					CreatedAt: parseDateTime("2024-01-01 10:00:00"),
					UpdatedAt: parseDateTime("2024-01-01 10:00:00"),
				}, nil)
			},
		},
		{
			name: "GetCustomerProfile Not Found",
			args: args{
				ctx: context.Background(),
				id:  2,
			},
			want:    nil,
			wantErr: true,
			mockFn: func(args args) {
				mockRepo.EXPECT().FindCustomerByID(gomock.Any(), args.id).Return(nil, errors.New(constants.ErrUserNotFound))
			},
		},
		{
			name: "GetCustomerProfile Internal Error",
			args: args{
				ctx: context.Background(),
				id:  3,
			},
			want:    nil,
			wantErr: true,
			mockFn: func(args args) {
				mockRepo.EXPECT().FindCustomerByID(gomock.Any(), args.id).Return(nil, errors.New("database connection failed"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)

			s := &customerService{
				customerRepository: mockRepo,
			}

			got, err := s.GetCustomerProfile(tt.args.ctx, tt.args.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("customerService.GetCustomerProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !assert.Equal(t, tt.want, got) {
				t.Errorf("customerService.GetCustomerProfile() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Helper functions for parsing dates and datetime
func parseDate(date string) time.Time {
	parsed, _ := time.Parse("2006-01-02", date)
	return parsed
}

func parseDateTime(datetime string) time.Time {
	parsed, _ := time.Parse("2006-01-02 15:04:05", datetime)
	return parsed
}
