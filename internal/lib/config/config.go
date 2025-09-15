package config

import (
	"log"
	"plata_currency_quotation/internal/lib/env"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

var V *Config

type Config struct {
	Env env.Environment `env:"ENV" env-required:"true"`

	QuotationUpdateIntervalMilliseconds int64 `env:"QUOTATION_UPDATE_INTERVAL_MILLISECONDS" env-required:"true"`

	DbHost     string `env:"DB_HOST" env-required:"true"`
	DbUser     string `env:"DB_USER" env-required:"true"`
	DbPassword string `env:"DB_PASSWORD" env-required:"true"`
	DbName     string `env:"DB_NAME" env-required:"true"`
	DbPort     uint16 `env:"DB_PORT" env-required:"true"`
	DbUseSsl   bool   `env:"DB_USE_SSL" env-required:"true"`

	ServerIp          string        `env:"SERVER_IP" env-required:"true"`
	ServerPort        uint16        `env:"SERVER_PORT" env-required:"true"`
	ServerTimeout     time.Duration `env:"SERVER_TIMEOUT" env-required:"true"`
	ServerIdleTimeout time.Duration `env:"SERVER_IDLE_TIMEOUT" env-required:"true"`

	SwaggerUser     string `env:"SWAGGER_USER"`
	SwaggerPassword string `env:"SWAGGER_PASSWORD"`

	MetricsPort uint16 `env:"METRICS_PORT" env-required:"true"`

	FrankfurterApiUrl string `env:"FRANKFURTER_API_URL" env-required:"true"`
}

func FromEnv() *Config {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		log.Fatalf("cannot read config: %s internal", err)
	}

	switch cfg.Env {
	case env.Local, env.Dev, env.Preprod, env.Prod:
	default:
		log.Fatalf("invalid Env value: %s internal", cfg.Env)
	}

	switch cfg.Env {
	case env.Dev, env.Preprod:
		if cfg.SwaggerUser == "" || cfg.SwaggerPassword == "" {
			log.Fatalf("Swagger Login and Swagger Password must be set in %s environment", cfg.Env)
		}
	}

	return &cfg
}

func Default() *Config {
	return &Config{
		Env:                                 env.Local,
		QuotationUpdateIntervalMilliseconds: 1000,
		DbHost:                              "",
		DbUser:                              "",
		DbPassword:                          "",
		DbName:                              "",
		DbPort:                              0,
		DbUseSsl:                            false,
		ServerIp:                            "localhost",
		ServerPort:                          8080,
		ServerTimeout:                       time.Second * 10,
		ServerIdleTimeout:                   time.Second * 60,
		MetricsPort:                         9090,
		FrankfurterApiUrl:                   "",
	}
}
