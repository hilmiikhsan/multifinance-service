package config

import (
	"sync"

	"github.com/rs/zerolog/log"

	"github.com/hilmiikhsan/multifinance-service/pkg/config"
	"github.com/hilmiikhsan/multifinance-service/pkg/utils"
)

var (
	Envs *Config // Envs is global vars Config.
	once sync.Once
)

type Config struct {
	App struct {
		Name                    string `env:"APP_NAME" env-default:"multifinance-service"`
		Environtment            string `env:"APP_ENV" env-default:"development"`
		BaseURL                 string `env:"APP_BASE_URL" env-default:"http://localhost:9090"`
		Port                    string `env:"APP_PORT" env-default:"9090"`
		LogLevel                string `env:"APP_LOG_LEVEL" env-default:"debug"`
		LogFile                 string `env:"APP_LOG_FILE" env-default:"./logs/app.log"`
		LogFileWs               string `env:"APP_LOG_FILE_WS" env-default:"./logs/ws.log"`
		LocalStoragePublicPath  string `env:"LOCAL_STORAGE_PUBLIC_PATH" env-default:"./storage/public"`
		LocalStoragePrivatePath string `env:"LOCAL_STORAGE_PRIVATE_PATH" env-default:"./storage/private"`
	}
	DB struct {
		ConnectionTimeout int `env:"DB_CONN_TIMEOUT" env-default:"30" env-description:"database timeout in seconds"`
		MaxOpenCons       int `env:"DB_MAX_OPEN_CONS" env-default:"20" env-description:"database max open conn in seconds"`
		MaxIdleCons       int `env:"DB_MAX_IdLE_CONS" env-default:"20" env-description:"database max idle conn in seconds"`
		ConnMaxLifetime   int `env:"DB_CONN_MAX_LIFETIME" env-default:"0" env-description:"database conn max lifetime in seconds"`
	}
	Guard struct {
		JwtPrivateKey             string `env:"JWT_PRIVATE_KEY" env-default:""`
		JwtTokenExpiration        string `env:"JWT_TOKEN_EXPIRATION" env-default:"15m"`
		JwtRefreshTokenExpiration string `env:"JWT_REFRESH_TOKEN_EXPIRATION" env-default:"72h"`
	}
	MultifinanceMysql struct {
		Host     string `env:"MULTIFINANCE_MYSQL_HOST" env-default:"localhost"`
		Port     string `env:"MULTIFINANCE_MYSQL_PORT" env-default:"8889"`
		Username string `env:"MULTIFINANCE_MYSQL_USER" env-default:"root"`
		Password string `env:"MULTIFINANCE_MYSQL_PASSWORD" env-default:"password"`
		Database string `env:"MULTIFINANCE_MYSQL_DB" env-default:"multifinance"`
		SslMode  string `env:"MULTIFINANCE_MYSQL_SSL_MODE" env-default:"disable"`
	}
	RedisDB struct {
		Host     string `env:"MULTIFINANCE_REDIS_HOST" env-default:"redis"`
		Port     string `env:"MULTIFINANCE_REDIS_PORT" env-default:"6379"`
		Password string `env:"MULTIFINANCE_REDIS_PASSWORD" env-default:"password"`
		Database int    `env:"MULTIFINANCE_REDIS_DB" env-default:"0"`
	}
}

// Option is Configure type return func.
type Option = func(c *Configure) error

// Configure is the data struct.
type Configure struct {
	path     string
	filename string
}

// Configuration create instance.
func Configuration(opts ...Option) *Configure {
	c := &Configure{}

	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			panic(err)
		}
	}
	return c
}

// Initialize will create instance of Configure.
func (c *Configure) Initialize() {
	once.Do(func() {
		Envs = &Config{}
		if err := config.Load(config.Opts{
			Config:    Envs,
			Paths:     []string{c.path},
			Filenames: []string{c.filename},
		}); err != nil {
			log.Fatal().Err(err).Msg("get config error")
		}

		Envs.App.Name = utils.GetEnv("APP_NAME", Envs.App.Name)
		Envs.App.Port = utils.GetEnv("APP_PORT", Envs.App.Port)
		Envs.App.LogLevel = utils.GetEnv("APP_LOG_LEVEL", Envs.App.LogLevel)
		Envs.App.LogFile = utils.GetEnv("APP_LOG_FILE", Envs.App.LogFile)
		Envs.App.LogFileWs = utils.GetEnv("APP_LOG_FILE_WS", Envs.App.LogFileWs)
		Envs.App.LocalStoragePublicPath = utils.GetEnv("LOCAL_STORAGE_PUBLIC_PATH", Envs.App.LocalStoragePublicPath)
		Envs.App.LocalStoragePrivatePath = utils.GetEnv("LOCAL_STORAGE_PRIVATE_PATH", Envs.App.LocalStoragePrivatePath)
		Envs.DB.ConnectionTimeout = utils.GetIntEnv("DB_CONN_TIMEOUT", Envs.DB.ConnectionTimeout)
		Envs.DB.MaxOpenCons = utils.GetIntEnv("DB_MAX_OPEN_CONS", Envs.DB.MaxOpenCons)
		Envs.DB.MaxIdleCons = utils.GetIntEnv("DB_MAX_IdLE_CONS", Envs.DB.MaxIdleCons)
		Envs.DB.ConnMaxLifetime = utils.GetIntEnv("DB_CONN_MAX_LIFETIME", Envs.DB.ConnMaxLifetime)
		Envs.Guard.JwtPrivateKey = utils.GetEnv("JWT_PRIVATE_KEY", Envs.Guard.JwtPrivateKey)
		Envs.Guard.JwtTokenExpiration = utils.GetEnv("JWT_TOKEN_EXPIRATION", Envs.Guard.JwtTokenExpiration)
		Envs.Guard.JwtRefreshTokenExpiration = utils.GetEnv("JWT_REFRESH_TOKEN_EXPIRATION", Envs.Guard.JwtRefreshTokenExpiration)
		Envs.MultifinanceMysql.Host = utils.GetEnv("MULTIFINANCE_MYSQL_HOST", Envs.MultifinanceMysql.Host)
		Envs.MultifinanceMysql.Port = utils.GetEnv("MULTIFINANCE_MYSQL_PORT", Envs.MultifinanceMysql.Port)
		Envs.MultifinanceMysql.Username = utils.GetEnv("MULTIFINANCE_MYSQL_USER", Envs.MultifinanceMysql.Username)
		Envs.MultifinanceMysql.Password = utils.GetEnv("MULTIFINANCE_MYSQL_PASSWORD", Envs.MultifinanceMysql.Password)
		Envs.MultifinanceMysql.Database = utils.GetEnv("MULTIFINANCE_MYSQL_DB", Envs.MultifinanceMysql.Database)
		Envs.MultifinanceMysql.SslMode = utils.GetEnv("MULTIFINANCE_MYSQL_SSL_MODE", Envs.MultifinanceMysql.SslMode)
		Envs.RedisDB.Host = utils.GetEnv("MULTIFINANCE_REDIS_HOST", Envs.RedisDB.Host)
		Envs.RedisDB.Port = utils.GetEnv("MULTIFINANCE_REDIS_PORT", Envs.RedisDB.Port)
		Envs.RedisDB.Password = utils.GetEnv("MULTIFINANCE_REDIS_PASSWORD", Envs.RedisDB.Password)
		Envs.RedisDB.Database = utils.GetIntEnv("MULTIFINANCE_REDIS_DB", Envs.RedisDB.Database)
	})
}

// WithPath will assign to field path Configure.
func WithPath(path string) Option {
	return func(c *Configure) error {
		c.path = path
		return nil
	}
}

// WithFilename will assign to field name Configure.
func WithFilename(name string) Option {
	return func(c *Configure) error {
		c.filename = name
		return nil
	}
}
