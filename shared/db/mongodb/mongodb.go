package mongodb

import (
	"context"
	"fmt"
	"net"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	Host       string
	Port       string
	User       string
	Password   string
	Database   string
	ReplicaSet string
}

func New(ctx context.Context, cfg Config) (*mongo.Client, error) {
	connURI := fmt.Sprintf("mongodb://%s/?directConnection=true&tls=false", net.JoinHostPort(cfg.Host, cfg.Port))
	option := options.Client().ApplyURI(connURI)

	client, err := mongo.Connect(ctx, option)

	if err != nil {
		return nil, fmt.Errorf("connect mongodb: %w", err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("ping mongodb: %w", err)
	}

	return client, nil
}
