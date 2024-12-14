package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type constants struct {
	App struct {
		Name         string         `yaml:"name"`
		Port         int            `yaml:"port"`
		ReadTimeout  int            `yaml:"read_timeout"`
		WriteTimeout int            `yaml:"write_timeout"`
		Timezone     string         `yaml:"timezone"`
		Debug        bool           `yaml:"debug"`
		Env          string         `yaml:"env"`
		SecretKey    string         `yaml:"secret_key"`
		ExpireIn     *time.Duration `yaml:"expire_in"`
	} `yaml:"App"`

	DB struct {
		DsnMain string `yaml:"dsn_main" env:"DSN_MAIN"`
	}
}

func TestLoad(t *testing.T) {
	var cfg constants
	_ = LoadFilePathEnv(Opts{
		Config:    &cfg,
		Filenames: []string{"test.yaml"},
		Paths:     []string{".", "./config"},
	})

	assert.Equal(t, 8778, cfg.App.Port)
}

func TestConfigPathFail(t *testing.T) {
	var cfg constants
	err := LoadFilePathEnv(Opts{
		Config:    &cfg,
		Filenames: []string{"test.yaml"},
		Paths:     []string{"./config"},
	})

	assert.Error(t, err)
}

func TestConfigEnv(t *testing.T) {
	os.Setenv("DSN_MAIN", "host=localhost port=5999 user=root password=root123 dbname=dbroot sslmode=disable")
	defer os.Unsetenv("DSN_MAIN")

	var cfg constants
	err := Load(Opts{
		Config:    &cfg,
		Filenames: []string{"test.yaml"},
		Paths:     []string{".", "./config"},
	})

	assert.NoError(t, err)
	assert.Equal(t, "host=localhost port=5999 user=root password=root123 dbname=dbroot sslmode=disable", cfg.DB.DsnMain)
}

func TestConfigFunc(t *testing.T) {
	os.Setenv("DSN_MAIN", "host=localhost port=5999 user=root password=root123 dbname=dbroot sslmode=disable")
	defer os.Unsetenv("DSN_MAIN")

	var cfg constants
	err := Load(Opts{
		Config:    &cfg,
		Filenames: []string{"test.yaml"},
		Paths:     []string{".", "./config"},
	})

	assert.NoError(t, err)
	assert.Equal(t, "host=localhost port=5999 user=root password=root123 dbname=dbroot sslmode=disable", cfg.DB.DsnMain)
}
