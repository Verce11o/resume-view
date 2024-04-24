package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

type Config struct {
	Env    string `env:"env"`
	Server Server
	DB     DB
	Jaeger Jaeger
}

type Server struct {
	Port string `env:"SERVER_PORT" env-default:"3007"`
}

type DB struct {
	User     string `env:"POSTGRES_USER" env-default:"postgres"`
	Password string `env:"POSTGRES_PASSWORD" env-default:"vercello"`
	Host     string `env:"POSTGRES_HOST" env-default:"localhost"`
	Port     string `env:"POSTGRES_PORT" env-default:"5432"`
	Name     string `env:"POSTGRES_DB" env-default:"views"`
	SSLMode  string `env:"POSTGRES_SSL_MODE" env-default:"disable"`
}

type Jaeger struct {
	Endpoint string `env:"JAEGER_ENDPOINT" env-default:"localhost:4317"`
}

func Load() *Config {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)

	if err != nil {
		log.Fatalf("error while read config: %v", err)
	}

	return &cfg
}
