package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hilmiikhsan/multifinance-service/constants"
	"github.com/hilmiikhsan/multifinance-service/internal/adapter"
	redisRepository "github.com/hilmiikhsan/multifinance-service/internal/infrastructure/redis"
	"github.com/hilmiikhsan/multifinance-service/internal/middleware"
	"github.com/hilmiikhsan/multifinance-service/internal/module/auth/dto"
	"github.com/hilmiikhsan/multifinance-service/internal/module/auth/ports"
	"github.com/hilmiikhsan/multifinance-service/internal/module/auth/service"
	customerRepository "github.com/hilmiikhsan/multifinance-service/internal/module/customer/repository"
	"github.com/hilmiikhsan/multifinance-service/pkg/err_msg"
	jwtHandler "github.com/hilmiikhsan/multifinance-service/pkg/jwt_handler"
	"github.com/hilmiikhsan/multifinance-service/pkg/response"
	"github.com/rs/zerolog/log"
)

type authHandler struct {
	service    ports.AuthService
	middleware middleware.AuthMiddleware
}

func NewAuthHandler() *authHandler {
	var handler = new(authHandler)

	// redis
	redisRepository := redisRepository.NewRedisRepository(adapter.Adapters.MultifinanceRedis)

	// jwt
	jwt := jwtHandler.NewJWT(redisRepository)

	// middleware
	middlewareHandler := middleware.NewAuthMiddleware(jwt)

	// repository
	customerRepository := customerRepository.NewCustomerRepository(adapter.Adapters.MultifinanceMysql)

	// service
	authService := service.NewUserService(
		adapter.Adapters.MultifinanceMysql,
		customerRepository,
		redisRepository,
		jwt,
	)

	// handler
	handler.service = authService
	handler.middleware = *middlewareHandler

	return handler
}

func (h *authHandler) AuthRoute(router fiber.Router) {
	router.Post("/register", h.register)
	router.Post("/login", h.login)
	router.Post("/refresh-token", h.refreshToken)
	router.Post("/logout", h.middleware.AuthBearer, h.logout)
}

func (h *authHandler) register(c *fiber.Ctx) error {
	var (
		req        = new(dto.RegisterRequest)
		ctx        = c.Context()
		validators = adapter.Adapters.Validator
	)

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::register - Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := validators.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::register - Invalid request body")
		code, errs := err_msg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	res, err := h.service.Register(ctx, req)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("handler::register - Failed to register user")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(res, ""))
}

func (h *authHandler) login(c *fiber.Ctx) error {
	var (
		req        = new(dto.LoginRequest)
		ctx        = c.Context()
		validators = adapter.Adapters.Validator
	)

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::login - Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := validators.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::login - Invalid request body")
		code, errs := err_msg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	res, err := h.service.Login(ctx, req)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("handler::login - Failed to login user")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}

func (h *authHandler) refreshToken(c *fiber.Ctx) error {
	var (
		ctx         = c.Context()
		accessToken = c.Get(constants.HeaderAuthorization)
	)

	if accessToken == "" {
		log.Warn().Msg("handler::refreshToken - Access token is required")
		return c.Status(fiber.StatusUnauthorized).JSON(response.Error(constants.ErrAccessTokenIsRequired))
	}

	if len(accessToken) > 7 {
		accessToken = accessToken[7:]
	}

	res, err := h.service.RefreshToken(ctx, accessToken)
	if err != nil {
		log.Error().Err(err).Any("access_token", accessToken).Msg("handler::refreshToken - Failed to refresh token")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}

func (h *authHandler) logout(c *fiber.Ctx) error {
	var (
		ctx         = c.Context()
		accessToken = c.Get(constants.HeaderAuthorization)
		locals      = middleware.GetLocals(c)
	)

	if accessToken == "" {
		log.Warn().Msg("handler::logout - Access token is required")
		return c.Status(fiber.StatusUnauthorized).JSON(response.Error(constants.ErrAccessTokenIsRequired))
	}

	if len(accessToken) > 7 {
		accessToken = accessToken[7:]
	}

	err := h.service.Logout(ctx, accessToken, locals)
	if err != nil {
		log.Error().Err(err).Any("access_token", accessToken).Msg("handler::logout - Failed to logout user")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(nil, ""))
}
