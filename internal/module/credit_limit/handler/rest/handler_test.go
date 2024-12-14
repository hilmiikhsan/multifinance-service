package rest

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/hilmiikhsan/multifinance-service/constants"
	"github.com/hilmiikhsan/multifinance-service/internal/module/credit_limit/dto"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func Test_creditLimitHandler_getCreditLimits(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockSvc := NewMockCreditLimitService(ctrlMock)

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
					mockSvc.EXPECT().GetCreditLimits(gomock.Any(), 1).Return(&[]dto.GetCreditLimitsResponse{
						{
							Tenor:       12,
							LimitAmount: 50000,
						},
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
				status: http.StatusInternalServerError,
				mockFn: func() {
					mockSvc.EXPECT().GetCreditLimits(gomock.Any(), 1).Return(nil, errors.New(constants.ErrInternalServerError))
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
					mockSvc.EXPECT().GetCreditLimits(gomock.Any(), 1).Return(nil, errors.New("internal server error"))
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			handler := &creditLimitHandler{service: mockSvc}
			app.Get("/limits", func(c *fiber.Ctx) error {
				if tt.args.token == "Bearer valid-token" {
					c.Locals("customer_id", 1)
				} else {
					return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
						"message": constants.ErrTokenAlreadyExpired,
						"success": false,
					})
				}
				return handler.getCreditLimits(c)
			})

			tt.args.mockFn()

			req := httptest.NewRequest(http.MethodGet, "/limits", nil)
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
