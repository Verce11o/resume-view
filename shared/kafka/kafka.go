package kafka

import (
	"context"
	"fmt"
	"net"

	"github.com/segmentio/kafka-go"
)

type Config struct {
	Host      string
	Port      string
	Topic     string
	Partition int
}

func New(ctx context.Context, cfg Config) (*kafka.Conn, error) {
	conn, err := kafka.DialLeader(ctx, "tcp", net.JoinHostPort(cfg.Host, cfg.Port), cfg.Topic, cfg.Partition)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to kafka: %w", err)
	}

	return conn, nil
}
