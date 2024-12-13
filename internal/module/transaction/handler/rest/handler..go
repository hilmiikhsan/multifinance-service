package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hilmiikhsan/multifinance-service/internal/adapter"
	redisRepository "github.com/hilmiikhsan/multifinance-service/internal/infrastructure/redis"
	"github.com/hilmiikhsan/multifinance-service/internal/middleware"
	creditLimitRepository "github.com/hilmiikhsan/multifinance-service/internal/module/credit_limit/repository"
	"github.com/hilmiikhsan/multifinance-service/internal/module/transaction/dto"
	"github.com/hilmiikhsan/multifinance-service/internal/module/transaction/ports"
	transactionRepository "github.com/hilmiikhsan/multifinance-service/internal/module/transaction/repository"
	"github.com/hilmiikhsan/multifinance-service/internal/module/transaction/service"
	"github.com/hilmiikhsan/multifinance-service/pkg/err_msg"
	jwtHandler "github.com/hilmiikhsan/multifinance-service/pkg/jwt_handler"
	"github.com/hilmiikhsan/multifinance-service/pkg/response"
	"github.com/rs/zerolog/log"
)

type transactionHandler struct {
	service    ports.TransactionService
	middleware middleware.AuthMiddleware
}

func NewTransactionHandler() *transactionHandler {
	var handler = new(transactionHandler)

	// redis
	redisRepository := redisRepository.NewRedisRepository(adapter.Adapters.MultifinanceRedis)

	// jwt
	jwt := jwtHandler.NewJWT(redisRepository)

	// middleware
	middlewareHandler := middleware.NewAuthMiddleware(jwt)

	// repository
	transactionRepository := transactionRepository.NewTransactionRepository(adapter.Adapters.MultifinanceMysql)
	creditLimitRepository := creditLimitRepository.NewCreditLimitRepository(adapter.Adapters.MultifinanceMysql)

	// service
	transactionService := service.NewTransactionService(
		adapter.Adapters.MultifinanceMysql,
		transactionRepository,
		creditLimitRepository,
	)

	// handler
	handler.service = transactionService
	handler.middleware = *middlewareHandler

	return handler
}

func (h *transactionHandler) TransactionRoute(router fiber.Router) {
	router.Post("/create", h.middleware.AuthBearer, h.createTranscation)
}

func (h *transactionHandler) createTranscation(c *fiber.Ctx) error {
	var (
		ctx        = c.Context()
		req        = new(dto.CreateTransactionRequest)
		locals     = middleware.GetLocals(c)
		validators = adapter.Adapters.Validator
	)

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::createTranscation - Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := validators.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::createTranscation - Invalid request body")
		code, errs := err_msg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	req.CustomerID = locals.GetCustomerID()

	if err := h.service.CreateTransaction(ctx, req); err != nil {
		log.Error().Err(err).Any("payload", req).Msg("handler::createTranscation - Failed to create transaction")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(nil, ""))
}
