package utils

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/hilmiikhsan/multifinance-service/constants"
	"github.com/hilmiikhsan/multifinance-service/pkg/err_msg"
	"github.com/rs/zerolog/log"
)

func HandleInsertUniqueError(err error, data interface{}, uniqueConstraints map[string]string) (interface{}, error) {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 { // 1062: Duplicate entry error
		errorMsg := mysqlErr.Message
		for key, customMessage := range uniqueConstraints {
			if strings.Contains(errorMsg, fmt.Sprintf("for key '%s'", key)) {
				log.Warn().Err(err).Any("payload", data).Msgf("repository::Insert - %s", customMessage)
				return nil, err_msg.NewCustomErrors(fiber.StatusConflict, err_msg.WithMessage(customMessage))
			}
		}

		// If no matching constraint was found
		log.Error().Err(err).Any("payload", data).Msg("repository::Insert - Unknown unique constraint violation")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// If it's not a MySQL duplicate entry error, log and return the original error
	log.Error().Err(err).Any("payload", data).Msg("repository::Insert - Failed to insert data")
	return nil, err
}
