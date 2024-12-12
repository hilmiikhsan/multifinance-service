package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hilmiikhsan/multifinance-service/internal/adapter"
	redisRepository "github.com/hilmiikhsan/multifinance-service/internal/infrastructure/redis"
	"github.com/hilmiikhsan/multifinance-service/internal/module/auth/dto"
	"github.com/hilmiikhsan/multifinance-service/internal/module/auth/ports"
	"github.com/hilmiikhsan/multifinance-service/internal/module/auth/service"
	customerRepository "github.com/hilmiikhsan/multifinance-service/internal/module/customer/repository"
	"github.com/hilmiikhsan/multifinance-service/pkg/err_msg"
	"github.com/hilmiikhsan/multifinance-service/pkg/response"
	"github.com/rs/zerolog/log"
)

type authHandler struct {
	service ports.AuthService
}

func NewAuthHandler() *authHandler {
	var handler = new(authHandler)

	// redis
	redisRepository := redisRepository.NewRedisRepository(adapter.Adapters.MultifinanceRedis)

	// repository
	customerRepository := customerRepository.NewCustomerRepository(adapter.Adapters.MultifinanceMysql)

	// service
	userService := service.NewUserService(
		adapter.Adapters.MultifinanceMysql,
		customerRepository,
		redisRepository,
	)

	// handler
	handler.service = userService

	return handler
}

func (h *authHandler) AuthRoute(router fiber.Router) {
	router.Post("/register", h.register)
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
