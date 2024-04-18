package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

type Config struct {
	Env               string `yaml:"env"`
	ClientTimeout     string `yaml:"client_timeout"`
	RetriesCount      string `yaml:"retries_count"`
	CompanyServiceURL string `yaml:"company_service_url"`
	ResumeServiceURL  string `yaml:"resume_service_url"`

	Server Server `yaml:"server"`
	DB     DB     `yaml:"db"`
	Jaeger Jaeger `yaml:"jaeger"`
	Nats   Nats   `yaml:"nats"`
}

type Server struct {
	Port string `yaml:"port"`
}

type DB struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Name     string `yaml:"name"`
	SSLMode  string `yaml:"ssl_mode"`
}

type Jaeger struct {
	Endpoint string `yaml:"endpoint"`
}

type Nats struct {
	Url string `yaml:"url" env:"NATS_URL"`
}

func Load() *Config {
	var cfg Config

	err := cleanenv.ReadConfig("config.yml", &cfg)

	if err != nil {
		log.Fatalf("error while read config: %v", err)
	}

	return &cfg
}
