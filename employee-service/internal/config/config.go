package config

import (
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HTTPServer    HTTPServer
	GRPCServer    GRPCServer
	Postgres      Postgres
	MongoDB       MongoDB
	Redis         Redis
	Kafka         Kafka
	JWTSignKey    string        `env:"JWT_SIGN_KEY" env-default:"jwt-sign-key"`
	TokenTTL      time.Duration `env:"TOKEN_TTL" env-default:"24h"`
	MainDatabase  string        `env:"MAIN_DATABASE" env-default:"postgres"`
	MainTransport string        `env:"MAIN_TRANSPORT" env-default:"http"`
	LogLevel      string        `env:"LOG_LEVEL" env-default:"DEBUG"`
}

type HTTPServer struct {
	Port   string `env:"HTTP_SERVER_PORT" env-default:":3009"`
	Router string `env:"HTTP_SERVER_ROUTER" env-default:"gorilla"`
}

type GRPCServer struct {
	Port string `env:"GRPC_SERVER_PORT" env-default:":3010"`
}

type Postgres struct {
	User     string `env:"POSTGRES_USER" env-default:"postgres"`
	Password string `env:"POSTGRES_PASSWORD" env-default:"vercello"`
	Host     string `env:"POSTGRES_HOST" env-default:"localhost"`
	Port     string `env:"POSTGRES_PORT" env-default:"5432"`
	Name     string `env:"POSTGRES_DB" env-default:"resume-views"`
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

type Kafka struct {
	Host  string `env:"KAFKA_HOST" env-default:"localhost"`
	Port  string `env:"KAFKA_PORT" env-default:"9092"`
	Topic string `env:"KAFKA_TOPIC" env-default:"employees-events"`
}

func Load() Config {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)

	if err != nil {
		log.Fatalf("error while read config: %v", err)
	}

	return cfg
}
