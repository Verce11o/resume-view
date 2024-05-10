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
		return nil, fmt.Errorf("position with id %s does not exist", positionID)
	}

	var position models.Position

	if err = json.Unmarshal(positionBytes, &position); err != nil {
		return nil, fmt.Errorf("failed to unmarshal position with id %s: %w", positionID, err)
	}

	return &position, nil
}

func (r *PositionCache) SetPosition(ctx context.Context, positionID string, position *models.Position) error {
	positionBytes, err := json.Marshal(position)

	if err != nil {
		return fmt.Errorf("failed to marshal position with id %s: %w", positionID, err)
	}

	err = r.client.Set(ctx, r.createKey(positionID), positionBytes, time.Second*time.Duration(positionTTL)).Err()

	if err != nil {
		return fmt.Errorf("failed to set position with id %s: %w", positionID, err)
	}

	return nil
}

func (r *PositionCache) DeletePosition(ctx context.Context, positionID string) error {
	err := r.client.Del(ctx, r.createKey(positionID)).Err()
	if err != nil {
		return fmt.Errorf("failed to delete position with id %s: %w", positionID, err)
	}

	return nil
}

func (r *PositionCache) createKey(key string) string {
	return fmt.Sprintf("position:%s", key)
}
