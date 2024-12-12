package adapter

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hilmiikhsan/multifinance-service/internal/infrastructure/config"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func WithMultifinanceMySQL() Option {
	return func(a *Adapter) {
		dbUser := config.Envs.MultifinanceMysql.Username
		dbPassword := config.Envs.MultifinanceMysql.Password
		dbName := config.Envs.MultifinanceMysql.Database
		dbHost := config.Envs.MultifinanceMysql.Host
		dbPort := config.Envs.MultifinanceMysql.Port

		dbMaxPoolSize := config.Envs.DB.MaxOpenCons
		dbMaxIdleConns := config.Envs.DB.MaxIdleCons
		dbConnMaxLifetime := config.Envs.DB.ConnMaxLifetime

		connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&loc=UTC",
			dbUser, dbPassword, dbHost, dbPort, dbName,
		)

		db, err := sqlx.Connect("mysql", connectionString)
		if err != nil {
			log.Fatal().Err(err).Msg("Error connecting to MySQL")
		}

		db.SetMaxOpenConns(dbMaxPoolSize)
		db.SetMaxIdleConns(dbMaxIdleConns)
		db.SetConnMaxLifetime(time.Duration(dbConnMaxLifetime) * time.Second)

		// Check connection
		err = db.Ping()
		if err != nil {
			log.Fatal().Err(err).Msg("Error connecting to Multifinance MySQL")
		}

		a.MultifinanceMysql = db
		log.Info().Msg("Multifinance MySQL connected")
	}
}
