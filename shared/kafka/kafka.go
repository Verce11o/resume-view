package kafka

import (
	"context"
	"fmt"
	"net"

	"github.com/segmentio/kafka-go"
)

type Config struct {
	Host string
	Port string
}

func New(ctx context.Context, cfg Config) (*kafka.Conn, error) {
	addr := net.JoinHostPort(cfg.Host, cfg.Port)

	conn, err := kafka.DialContext(ctx, "tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to kafka: %w", err)
	}

	return conn, nil
}
