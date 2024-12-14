package rest

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/hilmiikhsan/multifinance-service/constants"
	creditLimitDto "github.com/hilmiikhsan/multifinance-service/internal/module/credit_limit/dto"
	"github.com/hilmiikhsan/multifinance-service/internal/module/customer/dto"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_customerHandler_getCustomerProfile(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockSvc := NewMockCustomerService(ctrlMock)

	type args struct {
		token  string
		status int
		mockFn func()
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "JWT Valid - Success",
			args: args{
				token:  "Bearer valid-token",
				status: http.StatusOK,
				mockFn: func() {
					mockSvc.EXPECT().GetCustomerProfile(gomock.Any(), 1).Return(&dto.GetCustomerProfileResponse{
						ID:              1,
						Nik:             "123456789",
						FullName:        "Test User",
						LegalName:       "Test User",
						BirthPlace:      "City",
						BirthDate:       "1990-01-01",
						Salary:          10000,
						KtpPhotoPath:    "/path/to/ktp/photo",
						SelfiePhotoPath: "/path/to/selfie/photo",
						Limits: []creditLimitDto.CreditLimit{
							{Tenor: 12, LimitAmount: 50000},
							{Tenor: 24, LimitAmount: 100000},
						},
						CreatedAt: "2024-01-01 10:00:00",
						UpdatedAt: "2024-01-01 10:00:00",
					}, nil)
				},
			},
			wantErr: false,
		},
		{
			name: "JWT Invalid - Unauthorized",
			args: args{
				token:  "Bearer invalid-token",
				status: http.StatusUnauthorized,
				mockFn: func() {},
			},
			wantErr: true,
		},
		{
			name: "Profile Not Found",
			args: args{
				token:  "Bearer valid-token",
				status: http.StatusNotFound,
				mockFn: func() {
					mockSvc.EXPECT().GetCustomerProfile(gomock.Any(), 1).Return(nil, errors.New(constants.ErrUserNotFound))
				},
			},
			wantErr: true,
		},
		{
			name: "Internal Server Error",
			args: args{
				token:  "Bearer valid-token",
				status: http.StatusInternalServerError,
				mockFn: func() {
					mockSvc.EXPECT().GetCustomerProfile(gomock.Any(), 1).Return(nil, errors.New("internal server error"))
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			handler := &customerHandler{service: mockSvc}
			app.Get("/profile", func(c *fiber.Ctx) error {
				if tt.args.token == "Bearer valid-token" {
					c.Locals("customer_id", 1)
				} else {
					return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
						"message": constants.ErrTokenAlreadyExpired,
						"success": false,
					})
				}
				return handler.getCustomerProfile(c)
			})

			tt.args.mockFn()

			req := httptest.NewRequest(http.MethodGet, "/profile", nil)
			req.Header.Set("Authorization", tt.args.token)
			resp, _ := app.Test(req)

			assert.Equal(t, tt.args.status, resp.StatusCode)

			if resp.StatusCode == http.StatusOK {
				assert.NotNil(t, resp.Body)
			} else {
				assert.True(t, tt.wantErr)
			}
		})
	}
}
