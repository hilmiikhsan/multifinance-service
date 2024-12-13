package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hilmiikhsan/multifinance-service/internal/adapter"
	redisRepository "github.com/hilmiikhsan/multifinance-service/internal/infrastructure/redis"
	"github.com/hilmiikhsan/multifinance-service/internal/middleware"
	"github.com/hilmiikhsan/multifinance-service/internal/module/credit_limit/ports"
	creditLimitRepository "github.com/hilmiikhsan/multifinance-service/internal/module/credit_limit/repository"
	"github.com/hilmiikhsan/multifinance-service/internal/module/credit_limit/service"
	"github.com/hilmiikhsan/multifinance-service/pkg/err_msg"
	jwtHandler "github.com/hilmiikhsan/multifinance-service/pkg/jwt_handler"
	"github.com/hilmiikhsan/multifinance-service/pkg/response"
	"github.com/rs/zerolog/log"
)

type creditLimitHandler struct {
	service    ports.CreditLimitService
	middleware middleware.AuthMiddleware
}

func NewCreditLimitHandler() *creditLimitHandler {
	var handler = new(creditLimitHandler)

	// redis
	redisRepository := redisRepository.NewRedisRepository(adapter.Adapters.MultifinanceRedis)

	// jwt
	jwt := jwtHandler.NewJWT(redisRepository)

	// middleware
	middlewareHandler := middleware.NewAuthMiddleware(jwt)

	// repository
	creditLimitRepository := creditLimitRepository.NewCreditLimitRepository(adapter.Adapters.MultifinanceMysql)

	// service
	creditLimitervice := service.NewCreditLimitService(
		adapter.Adapters.MultifinanceMysql,
		creditLimitRepository,
	)

	// handler
	handler.service = creditLimitervice
	handler.middleware = *middlewareHandler

	return handler
}

func (h *creditLimitHandler) CreditLimitRoute(router fiber.Router) {
	router.Get("/limits", h.middleware.AuthBearer, h.getCreditLimits)
}

func (h *creditLimitHandler) getCreditLimits(c *fiber.Ctx) error {
	var (
		ctx    = c.Context()
		locals = middleware.GetLocals(c)
	)

	res, err := h.service.GetCreditLimits(ctx, locals.GetCustomerID())
	if err != nil {
		log.Error().Err(err).Any("payload", locals.CustomerID).Msg("handler::CreditLimitRoute - Failed to get credit limit user")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}
