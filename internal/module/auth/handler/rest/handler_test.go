package rest

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/hilmiikhsan/multifinance-service/internal/module/auth/dto"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_authHandler_register(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockSvc := NewMockAuthService(ctrlMock)
	mockValidator := NewMockValidator(ctrlMock)

	type args struct {
		body       string
		statusCode int
		mockFn     func()
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Success - User Registered",
			args: args{
				body: `{
					"username": "testuser",
					"password": "password123",
					"email": "test@example.com"
				}`,
				statusCode: http.StatusCreated,
				mockFn: func() {
					mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
					mockSvc.EXPECT().Register(gomock.Any(), gomock.Any()).Return(&dto.RegisterResponse{
						ID:    1,
						Email: "test@example.com",
					}, nil)
				},
			},
			wantErr: false,
		},
		{
			name: "Failure - Invalid JSON Body",
			args: args{
				body:       `invalid-body`,
				statusCode: http.StatusBadRequest,
				mockFn:     func() {},
			},
			wantErr: true,
		},
		{
			name: "Failure - Validation Error",
			args: args{
				body: `{
					"username": "",
					"password": "short",
					"email": "invalid-email"
				}`,
				statusCode: http.StatusBadRequest,
				mockFn: func() {
					mockValidator.EXPECT().Validate(gomock.Any()).Return(errors.New("validation error"))
				},
			},
			wantErr: true,
		},
		{
			name: "Failure - Service Error",
			args: args{
				body: `{
					"username": "testuser",
					"password": "password123",
					"email": "test@example.com"
				}`,
				statusCode: http.StatusInternalServerError,
				mockFn: func() {
					mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
					mockSvc.EXPECT().Register(gomock.Any(), gomock.Any()).Return(nil, errors.New("service error"))
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			handler := &authHandler{
				service:   mockSvc,
				validator: mockValidator,
			}

			app.Post("/register", handler.register)

			tt.args.mockFn()

			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(tt.args.body))
			req.Header.Set("Content-Type", "application/json")
			resp, _ := app.Test(req)

			assert.Equal(t, tt.args.statusCode, resp.StatusCode)

			if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
				assert.NotNil(t, resp.Body)
			} else {
				assert.Equal(t, tt.args.statusCode, resp.StatusCode)
			}
		})
	}
}

func Test_authHandler_login(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockSvc := NewMockAuthService(ctrlMock)
	mockValidator := NewMockValidator(ctrlMock)

	type args struct {
		body       string
		statusCode int
		mockFn     func()
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Success - Valid login request",
			args: args{
				body:       `{"email":"test@example.com","password":"validpassword"}`,
				statusCode: fiber.StatusOK,
				mockFn: func() {
					mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
					mockSvc.EXPECT().Login(gomock.Any(), gomock.Any()).Return(&dto.LoginResponse{
						Token: "testtoken",
					}, nil)
				},
			},
			wantErr: false,
		},
		{
			name: "Failure - Invalid request body",
			args: args{
				body: `{
					"email": "invalid-email",
					"password": "short"
				}`,
				statusCode: http.StatusBadRequest,
				mockFn: func() {
					mockValidator.EXPECT().Validate(gomock.Any()).Return(errors.New("validation error"))
				},
			},
			wantErr: true,
		},
		{
			name: "Failure - Login service error",
			args: args{
				body:       `{"email":"test@example.com","password":"validpassword"}`,
				statusCode: fiber.StatusInternalServerError,
				mockFn: func() {
					mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
					mockSvc.EXPECT().Login(gomock.Any(), gomock.Any()).Return(nil, errors.New("invalid credentials"))
				},
			},
			wantErr: true,
		},
		{
			name: "Failure - Failed to parse request body",
			args: args{
				body:       `invalid-json-body`,
				statusCode: fiber.StatusBadRequest,
				mockFn: func() {
					// No validation or service call expected here
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			handler := &authHandler{
				service:   mockSvc,
				validator: mockValidator,
			}

			app.Post("/login", handler.login)

			tt.args.mockFn()

			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(tt.args.body))
			req.Header.Set("Content-Type", "application/json")
			resp, _ := app.Test(req)

			assert.Equal(t, tt.args.statusCode, resp.StatusCode)

			if resp.StatusCode == http.StatusOK {
				assert.NotNil(t, resp.Body)
			} else {
				assert.Equal(t, tt.args.statusCode, resp.StatusCode)
			}
		})
	}
}
