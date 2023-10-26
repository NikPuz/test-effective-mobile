package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	ShutdownTimeout time.Duration `mapstructure:"SHUTDOWN_TIMEOUT" validate:"required"`

	Port         int           `mapstructure:"SERVER_PORT" validate:"required"`
	WriteTimeout time.Duration `mapstructure:"SERVER_WRITE_TIMEOUT" validate:"required"`
	ReadTimeout  time.Duration `mapstructure:"SERVER_READ_TIMEOUT" validate:"required"`
	IdleTimeout  time.Duration `mapstructure:"SERVER_IDLE_TIMEOUT" validate:"required"`

	Username        string        `mapstructure:"DB_USERNAME" validate:"required"`
	Password        string        `mapstructure:"DB_PASSWORD" validate:"required"`
	Address         string        `mapstructure:"DB_ADDRESS"`
	DBName          string        `mapstructure:"DB_NAME" validate:"required"`
	Params          string        `mapstructure:"DB_PARAMS"`
	MaxConnLifetime time.Duration `mapstructure:"DB_MAX_CONN_LIFETIME" validate:"required"`
	MaxConnIdleTime time.Duration `mapstructure:"DB_MAX_CONN_IDLE_TIME" validate:"required"`
	MaxOpenCons     int           `mapstructure:"DB_MAX_OPEN_CONS" validate:"required"`
	MaxIdleCons     int           `mapstructure:"DB_MAX_IDLE_CONS" validate:"required"`

	EnrichmentAgeDomain         string        `mapstructure:"ENRICHMENT_AGE_DOMAIN" validate:"required"`
	EnrichmentGenderDomain      string        `mapstructure:"ENRICHMENT_GENDER_DOMAIN" validate:"required"`
	EnrichmentNationalityDomain string        `mapstructure:"ENRICHMENT_NATIONALITY_DOMAIN" validate:"required"`
	EnrichmentTimeout           time.Duration `mapstructure:"ENRICHMENT_TIMEOUT" validate:"required"`
}

func NewConfig() *Config {
	cfg := new(Config)

	viper.SetConfigType("env")
	viper.AddConfigPath("./configs/")
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		panic(err.Error())
	}
	if err := viper.Unmarshal(cfg); err != nil {
		panic(err.Error())
	}
	if err := validator.New().Struct(cfg); err != nil {
		panic(err.Error())
	}

	return cfg
}
