package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Verce11o/resume-view/employee-service/internal/models"
	"github.com/redis/go-redis/v9"
)

const (
	positionTTL = 3600
)

type PositionCache struct {
	client *redis.Client
}

func NewPositionCache(client *redis.Client) *PositionCache {
	return &PositionCache{client: client}
}

func (r *PositionCache) GetPosition(ctx context.Context, positionID string) (*models.Position, error) {
	positionBytes, err := r.client.Get(ctx, r.createKey(positionID)).Bytes()

	if err != nil || errors.Is(err, redis.Nil) {
		return nil, err
	}

	var position models.Position

	if err = json.Unmarshal(positionBytes, &position); err != nil {
		return nil, err
	}

	return &position, nil
}

func (r *PositionCache) SetPosition(ctx context.Context, positionID string, position *models.Position) error {
	positionBytes, err := json.Marshal(position)

	if err != nil {
		return err
	}

	return r.client.Set(ctx, r.createKey(positionID), positionBytes, time.Second*time.Duration(positionTTL)).Err()
}

func (r *PositionCache) DeletePosition(ctx context.Context, positionID string) error {
	return r.client.Del(ctx, r.createKey(positionID)).Err()
}

func (r *PositionCache) createKey(key string) string {
	return fmt.Sprintf("position:%s", key)
}
