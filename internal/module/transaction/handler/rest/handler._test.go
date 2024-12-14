package rest

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/hilmiikhsan/multifinance-service/internal/middleware"
	"github.com/hilmiikhsan/multifinance-service/internal/module/transaction/dto"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func Test_transactionHandler_createTranscation(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockSvc := NewMockTransactionService(ctrlMock)
	mockValidator := NewMockValidator(ctrlMock)

	type args struct {
		body   string
		mockFn func(*MockValidator, *MockTransactionService)
	}

	tests := []struct {
		name           string
		args           args
		expectedStatus int
		wantErr        bool
	}{
		{
			name: "Success - Create Transaction",
			args: args{
				body: `{
					"on_the_road_price": 500000,
					"installment_amount": 50000,
					"interest_amount": 5000,
					"asset_name": "Yamaha NMAX",
					"tenor_month": 12
				}`,
				mockFn: func(mv *MockValidator, ms *MockTransactionService) {
					mv.EXPECT().Validate(gomock.Any()).Return(nil)
					ms.EXPECT().CreateTransaction(gomock.Any(), gomock.Any()).Return(nil)
				},
			},
			expectedStatus: fiber.StatusCreated,
			wantErr:        false,
		},
		{
			name: "Failure - Invalid JSON Body",
			args: args{
				body: `invalid-json-body`,
				mockFn: func(mv *MockValidator, ms *MockTransactionService) {
					// No expectations for this test case
				},
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name: "Failure - Validation Error",
			args: args{
				body: `{
					"on_the_road_price": 0,
					"installment_amount": 0,
					"interest_amount": 0,
					"asset_name": "",
					"tenor_month": 0
				}`,
				mockFn: func(mv *MockValidator, ms *MockTransactionService) {
					mv.EXPECT().Validate(gomock.Any()).Return(errors.New("validation error"))
				},
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name: "Failure - Service Error",
			args: args{
				body: `{
					"on_the_road_price": 500000,
					"installment_amount": 50000,
					"interest_amount": 5000,
					"asset_name": "Yamaha NMAX",
					"tenor_month": 12
				}`,
				mockFn: func(mv *MockValidator, ms *MockTransactionService) {
					mv.EXPECT().Validate(gomock.Any()).Return(nil)
					ms.EXPECT().CreateTransaction(gomock.Any(), gomock.Any()).Return(errors.New("service error"))
				},
			},
			expectedStatus: fiber.StatusInternalServerError,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()

			mockLocals := &middleware.Locals{
				CustomerID: 1,
			}

			handler := &transactionHandler{
				service:   mockSvc,
				validator: mockValidator,
			}

			app.Post("/create", func(c *fiber.Ctx) error {
				c.Locals(mockLocals)
				return handler.createTranscation(c)
			})

			tt.args.mockFn(mockValidator, mockSvc)

			req := httptest.NewRequest(http.MethodPost, "/create", bytes.NewBufferString(tt.args.body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode, "Unexpected status code")

			if tt.expectedStatus == fiber.StatusCreated {
				assert.NotNil(t, resp.Body)
			}
		})
	}
}

func Test_transactionHandler_getDetailTransaction(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockSvc := NewMockTransactionService(ctrlMock)

	type args struct {
		id     string
		mockFn func(*MockTransactionService)
	}

	tests := []struct {
		name           string
		args           args
		expectedStatus int
		wantErr        bool
		customerID     int
	}{
		{
			name: "JWT Valid - Success",
			args: args{
				id: "1",
				mockFn: func(ms *MockTransactionService) {
					ms.EXPECT().GetDetailTransaction(
						gomock.Any(),
						1,
						1,
					).Return(&dto.GetDetailTransactionResponse{
						ID:                1,
						CustomerID:        1,
						ContractNumber:    "123456",
						OnTheRoadPrice:    500000,
						AdminFee:          5000,
						InstallmentAmount: 50000,
						InterestAmount:    5000,
					}, nil)
				},
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
			customerID:     1,
		},
		{
			name: "ID is Zero",
			args: args{
				id:     "0",
				mockFn: func(ms *MockTransactionService) {},
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
			customerID:     1,
		},
		{
			name: "Invalid ID Format",
			args: args{
				id:     ":id",
				mockFn: func(ms *MockTransactionService) {},
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
			customerID:     1,
		},
		{
			name: "Internal Server Error",
			args: args{
				id: "1",
				mockFn: func(ms *MockTransactionService) {
					ms.EXPECT().GetDetailTransaction(
						gomock.Any(),
						1,
						1,
					).Return(nil, errors.New("internal server error"))
				},
			},
			expectedStatus: fiber.StatusInternalServerError,
			wantErr:        true,
			customerID:     1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			handler := &transactionHandler{service: mockSvc}

			app.Get("/:id", func(c *fiber.Ctx) error {
				c.Locals("customer_id", tt.customerID)
				return handler.getDetailTransaction(c)
			})

			tt.args.mockFn(mockSvc)

			req := httptest.NewRequest(http.MethodGet, "/"+tt.args.id, nil)
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode, "Unexpected status code")

			if tt.expectedStatus == fiber.StatusOK {
				assert.NotNil(t, resp.Body)
			} else {
				assert.True(t, tt.wantErr)
			}
		})
	}
}

func Test_transactionHandler_getHistoryListTransaction(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockSvc := NewMockTransactionService(ctrlMock)
	mockValidator := NewMockValidator(ctrlMock)

	type args struct {
		id     string
		mockFn func(*MockTransactionService)
	}

	tests := []struct {
		name           string
		args           args
		expectedStatus int
		wantErr        bool
	}{
		{
			name: "JWT Valid - Success",
			args: args{
				id: "1",
				mockFn: func(ms *MockTransactionService) {
					ms.EXPECT().GetDetailTransaction(
						gomock.Any(),
						1,
						1,
					).Return(&dto.GetDetailTransactionResponse{
						ID:                1,
						CustomerID:        1,
						ContractNumber:    "123456",
						OnTheRoadPrice:    500000,
						AdminFee:          5000,
						InstallmentAmount: 50000,
						InterestAmount:    5000,
					}, nil)
				},
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		{
			name: "ID is Zero",
			args: args{
				id:     "0",
				mockFn: func(ms *MockTransactionService) {},
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name: "Invalid ID Format",
			args: args{
				id:     ":id",
				mockFn: func(ms *MockTransactionService) {},
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name: "Internal Server Error",
			args: args{
				id: "1",
				mockFn: func(ms *MockTransactionService) {
					ms.EXPECT().GetDetailTransaction(
						gomock.Any(),
						1,
						1,
					).Return(nil, errors.New("internal server error"))
				},
			},
			expectedStatus: fiber.StatusInternalServerError,
			wantErr:        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			handler := &transactionHandler{
				service:   mockSvc,
				validator: mockValidator,
			}

			app.Get("/:id", func(c *fiber.Ctx) error {
				c.Locals("customer_id", 1)
				return handler.getDetailTransaction(c)
			})

			tt.args.mockFn(mockSvc)

			req := httptest.NewRequest(http.MethodGet, "/"+tt.args.id, nil)
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode, "Unexpected status code")

			if tt.expectedStatus == fiber.StatusOK {
				assert.NotNil(t, resp.Body)
			} else {
				assert.True(t, tt.wantErr)
			}
		})
	}
}
