package adapter

import (
	"errors"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var (
	Adapters *Adapter
)

type Option func(adapter *Adapter)

//go:generate mockgen -source=adapters.go -destination=service_validator_mock_test.go -package=adapter
type Validator interface {
	Validate(i any) error
}

type Adapter struct {
	// Driving Adapters
	RestServer *fiber.App

	//Driven Adapters
	MultifinanceMysql *sqlx.DB
	MultifinanceRedis *redis.Client
	Validator         Validator // *validator.Validator
}

func (a *Adapter) Sync(opts ...Option) error {
	var errs []string

	for _, opt := range opts {
		opt(a)
	}

	if a.MultifinanceMysql == nil {
		errs = append(errs, "Multifinance Mysql not initialized")
	}

	if a.MultifinanceRedis == nil {
		errs = append(errs, "Multifinance Redis not initialized")
	}

	if a.RestServer == nil {
		errs = append(errs, "No server initialized")
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}

	return nil
}

func (a *Adapter) Unsync() error {
	var errs []string

	if a.RestServer != nil {
		if err := a.RestServer.Shutdown(); err != nil {
			errs = append(errs, err.Error())
		}
		log.Info().Msg("Rest server disconnected")
	}

	if a.MultifinanceMysql != nil {
		if err := a.MultifinanceMysql.Close(); err != nil {
			errs = append(errs, err.Error())
		}
		log.Info().Msg("Multifinance Mysql disconnected")
	}

	if a.MultifinanceRedis != nil {
		if err := a.MultifinanceRedis.Close(); err != nil {
			errs = append(errs, err.Error())
		}
		log.Info().Msg("Multifinance Redis disconnected")
	}

	if len(errs) > 0 {
		err := errors.New(strings.Join(errs, "\n"))
		log.Error().Msgf("Error while disconnecting adapters: %v", err)
		return err
	}

	return nil
}
