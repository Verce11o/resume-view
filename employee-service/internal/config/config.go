package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

type Config struct {
	Server       Server
	Postgres     Postgres
	MainDatabase string `env:"MAIN_DATABASE" envDefault:"postgres"`
	LogLevel     string `env:"LOG_LEVEL" env-default:"INFO"`
}

type Server struct {
	Port string `env:"SERVER_PORT" env-default:":3009"`
}

type Postgres struct {
	User     string `env:"POSTGRES_USER" env-default:"postgres"`
	Password string `env:"POSTGRES_PASSWORD" env-default:"vercello"`
	Host     string `env:"POSTGRES_HOST" env-default:"localhost"`
	Port     string `env:"POSTGRES_PORT" env-default:"5432"`
	Name     string `env:"POSTGRES_DB" env-default:"employees"`
	SSLMode  string `env:"POSTGRES_SSL_MODE" env-default:"disable"`
}

func Load() Config {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)

	if err != nil {
		log.Fatalf("error while read config: %v", err)
	}

	return cfg
}
