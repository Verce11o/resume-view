package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	LogLevel   string `env:"LOG_LEVEL" env-default:"DEBUG"`
	GRPCServer GRPCServer
	HTTPServer HTTPServer
	DB         DB
	Kafka      Kafka
	Jaeger     Jaeger
}

type GRPCServer struct {
	Port string `env:"GRPC_SERVER_PORT" env-default:"3007"`
}

type HTTPServer struct {
	Port string `env:"HTTP_SERVER_PORT" env-default:":3030"`
}

type DB struct {
	User     string `env:"POSTGRES_USER" env-default:"postgres"`
	Password string `env:"POSTGRES_PASSWORD" env-default:"vercello"`
	Host     string `env:"POSTGRES_HOST" env-default:"localhost"`
	Port     string `env:"POSTGRES_PORT" env-default:"5432"`
	Name     string `env:"POSTGRES_DB" env-default:"views"`
	SSLMode  string `env:"POSTGRES_SSL_MODE" env-default:"disable"`
}

type Kafka struct {
	Host    string `env:"KAFKA_HOST" env-default:"localhost"`
	Port    string `env:"KAFKA_PORT" env-default:"9092"`
	Topic   string `env:"KAFKA_TOPIC" env-default:"employees-events"`
	GroupID string `env:"KAFKA_GROUP_ID" env-default:"Group1"`
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
