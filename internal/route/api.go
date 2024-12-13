package route

import (
	"github.com/hilmiikhsan/multifinance-service/pkg/response"

	"github.com/gofiber/fiber/v2"
	authRest "github.com/hilmiikhsan/multifinance-service/internal/module/auth/handler/rest"
	creditLimitRest "github.com/hilmiikhsan/multifinance-service/internal/module/credit_limit/handler/rest"
	customerRest "github.com/hilmiikhsan/multifinance-service/internal/module/customer/handler/rest"
	"github.com/rs/zerolog/log"
)

func SetupRoutes(app *fiber.App) {
	var (
		authAPIV1        = app.Group("/api/v1/auth")
		customerAPIV1    = app.Group("/api/v1/customer")
		creditLimitAPIV1 = app.Group("/api/v1/credit")
	)

	authRest.NewAuthHandler().AuthRoute(authAPIV1)
	customerRest.NewCustomerHandler().CustomerRoute(customerAPIV1)
	creditLimitRest.NewCreditLimitHandler().CreditLimitRoute(creditLimitAPIV1)

	// fallback route
	app.Use(func(c *fiber.Ctx) error {
		var (
			method = c.Method()                       // get the request method
			path   = c.Path()                         // get the request path
			query  = c.Context().QueryArgs().String() // get all query params
			ua     = c.Get("User-Agent")              // get the request user agent
			ip     = c.IP()                           // get the request IP
		)

		log.Info().
			Str("url", c.OriginalURL()).
			Str("method", method).
			Str("path", path).
			Str("query", query).
			Str("ua", ua).
			Str("ip", ip).
			Msg("Route not found.")
		return c.Status(fiber.StatusNotFound).JSON(response.Error("Route not found"))
	})
}
