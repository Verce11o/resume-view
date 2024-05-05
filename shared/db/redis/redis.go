package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

type Config struct {
	Host     string
	Port     string
	Password string
	Database int
}

func New(ctx context.Context, cfg Config) (*redis.Client, error) {
	redisHost := fmt.Sprintf("%v:%v", cfg.Host, cfg.Port)

	client := redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: cfg.Password,
		DB:       cfg.Database,
	})

	_, err := client.Ping(ctx).Result()

	if err != nil {
		return nil, err
	}

	return client, nil
}
