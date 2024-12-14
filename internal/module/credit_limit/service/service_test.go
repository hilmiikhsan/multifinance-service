package service

import (
	"context"
	"errors"
	"testing"

	"github.com/hilmiikhsan/multifinance-service/internal/module/credit_limit/dto"
	entity "github.com/hilmiikhsan/multifinance-service/internal/module/credit_limit/entity"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func Test_creditLimitService_GetCreditLimits(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockRepo := NewMockCreditLimitRepository(ctrlMock)

	type args struct {
		ctx        context.Context
		customerID int
	}
	tests := []struct {
		name    string
		args    args
		want    *[]dto.GetCreditLimitsResponse
		wantErr bool
		mockFn  func(args args)
	}{
		{
			name: "GetCreditLimits Success",
			args: args{
				ctx:        context.Background(),
				customerID: 1,
			},
			want: &[]dto.GetCreditLimitsResponse{
				{
					Tenor:       12,
					LimitAmount: 50000,
				},
				{
					Tenor:       24,
					LimitAmount: 100000,
				},
			},
			wantErr: false,
			mockFn: func(args args) {
				mockRepo.EXPECT().FindCreditLimitByCustomerID(gomock.Any(), args.customerID).Return(&[]entity.Limits{
					{
						TenorMonth:  12,
						LimitAmount: 50000,
					},
					{
						TenorMonth:  24,
						LimitAmount: 100000,
					},
				}, nil)
			},
		},
		{
			name: "GetCreditLimits Internal Error",
			args: args{
				ctx:        context.Background(),
				customerID: 3,
			},
			want:    nil,
			wantErr: true,
			mockFn: func(args args) {
				mockRepo.EXPECT().FindCreditLimitByCustomerID(gomock.Any(), args.customerID).Return(nil, errors.New("database connection failed"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)

			s := &creditLimitService{
				creditLimitRepository: mockRepo,
			}
			got, err := s.GetCreditLimits(tt.args.ctx, tt.args.customerID)

			if tt.wantErr {
				assert.Error(t, err, "expected an error but got none")
			} else {
				assert.NoError(t, err, "did not expect an error but got one")
			}

			assert.Equal(t, tt.want, got, "unexpected result")
		})
	}
}
