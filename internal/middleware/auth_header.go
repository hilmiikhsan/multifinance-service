package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hilmiikhsan/multifinance-service/constants"
	"github.com/rs/zerolog/log"
)

func (m *AuthMiddleware) AuthBearer(c *fiber.Ctx) error {
	accessToken := c.Get(constants.HeaderAuthorization)
	unauthorizedResponse := fiber.Map{
		"message": "Unauthorized",
		"success": false,
	}

	// If the cookie is not set, return an unauthorized status
	if accessToken == "" {
		log.Error().Msg("middleware::AuthBearer - Unauthorized [Header not set]")
		return c.Status(fiber.StatusUnauthorized).JSON(unauthorizedResponse)
	}

	// remove the Bearer prefix
	if len(accessToken) > 7 {
		accessToken = accessToken[7:]
	}

	// Parse the JWT string and store the result in `claims`
	claims, err := m.jwt.ParseTokenString(c.Context(), accessToken)
	if err != nil {
		log.Error().Err(err).Any("payload", accessToken).Msg("middleware::AuthBearer - Error while parsing token")
		return c.Status(fiber.StatusUnauthorized).JSON(unauthorizedResponse)
	}

	c.Locals("user_id", claims.UserId)
	c.Locals("nik", claims.Nik)
	c.Locals("email", claims.Email)
	c.Locals("full_name", claims.FullName)

	// If the token is valid, pass the request to the next handler
	return c.Next()
}
