package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hilmiikhsan/multifinance-service/internal/adapter"
	redisRepository "github.com/hilmiikhsan/multifinance-service/internal/infrastructure/redis"
	"github.com/hilmiikhsan/multifinance-service/internal/middleware"
	"github.com/hilmiikhsan/multifinance-service/internal/module/customer/ports"
	customerRepository "github.com/hilmiikhsan/multifinance-service/internal/module/customer/repository"
	"github.com/hilmiikhsan/multifinance-service/internal/module/customer/service"
	"github.com/hilmiikhsan/multifinance-service/pkg/err_msg"
	jwtHandler "github.com/hilmiikhsan/multifinance-service/pkg/jwt_handler"
	"github.com/hilmiikhsan/multifinance-service/pkg/response"
	"github.com/rs/zerolog/log"
)

type customerHandler struct {
	service    ports.CustomerService
	middleware middleware.AuthMiddleware
}

func NewCustomerHandler() *customerHandler {
	var handler = new(customerHandler)

	// redis
	redisRepository := redisRepository.NewRedisRepository(adapter.Adapters.MultifinanceRedis)

	// jwt
	jwt := jwtHandler.NewJWT(redisRepository)

	// middleware
	middlewareHandler := middleware.NewAuthMiddleware(jwt)

	// repository
	customerRepository := customerRepository.NewCustomerRepository(adapter.Adapters.MultifinanceMysql)

	// service
	customerService := service.NewCustomerService(
		adapter.Adapters.MultifinanceMysql,
		customerRepository,
	)

	// handler
	handler.service = customerService
	handler.middleware = *middlewareHandler

	return handler
}

func (h *customerHandler) CustomerRoute(router fiber.Router) {
	router.Get("/profile", h.middleware.AuthBearer, h.getCustomerProfile)
}

func (h *customerHandler) getCustomerProfile(c *fiber.Ctx) error {
	var (
		ctx    = c.Context()
		locals = middleware.GetLocals(c)
	)

	res, err := h.service.GetCustomerProfile(ctx, locals.GetCustomerID())
	if err != nil {
		log.Error().Err(err).Any("response", res).Msg("handler::getCustomerProfile - Failed to logout user")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}
