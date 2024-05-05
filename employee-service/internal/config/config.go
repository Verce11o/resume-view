package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

type Config struct {
	Server       Server
	Postgres     Postgres
	MongoDB      MongoDB
	Redis        Redis
	MainDatabase string `env:"MAIN_DATABASE" envDefault:"postgres"`
	LogLevel     string `env:"LOG_LEVEL" env-default:"DEBUG"`
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

type MongoDB struct {
	User       string `env:"MONGO_USER" env-default:""`
	Password   string `env:"MONGO_PASSWORD" env-default:""`
	Host       string `env:"MONGO_HOST" env-default:"localhost"`
	Port       string `env:"MONGO_PORT" env-default:"27017"`
	Name       string `env:"MONGO_DB" env-default:"employees"`
	ReplicaSet string `env:"MONGO_REPLICA_SET" env-default:"rs0"`
}

type Redis struct {
	Host     string `env:"REDIS_HOST" env-default:"localhost"`
	Port     string `env:"REDIS_PORT" env-default:"6379"`
	Password string `env:"REDIS_PASSWORD" env-default:""`
	Database int    `env:"REDIS_DB" env-default:"0"`
}

func Load() Config {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)

	if err != nil {
		log.Fatalf("error while read config: %v", err)
	}

	return cfg
}
