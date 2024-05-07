package redis

import (
	"context"
	"github.com/Verce11o/resume-view/employee-service/internal/models"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	redisContainer "github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"
	"testing"
	"time"
)

func setupRedisContainer(t *testing.T) (*redisContainer.RedisContainer, string) {
	ctx := context.Background()
	container, err := redisContainer.RunContainer(ctx,
		testcontainers.WithImage("redis:latest"),
		testcontainers.WithWaitStrategy(
			wait.
				ForLog("Ready to accept connections tcp").
				WithStartupTimeout(3*time.Second),
		),
	)
	require.NoError(t, err)

	connURI, err := container.Endpoint(ctx, "")
	require.NoError(t, err)

	return container, connURI
}

func TestPositionCache_SetPosition(t *testing.T) {

	ctx := context.Background()

	container, connURI := setupRedisContainer(t)
	defer func(container *redisContainer.RedisContainer, ctx context.Context) {
		err := container.Terminate(ctx)
		if err != nil {
			t.Fatal(err)
		}
	}(container, ctx)

	client := redis.NewClient(&redis.Options{
		Addr: connURI,
	})

	_, err := client.Ping(ctx).Result()

	require.NoError(t, err)

	positionID := uuid.New().String()
	positionCacheRepo := NewPositionCache(client)

	tests := []struct {
		name       string
		positionID string
		position   *models.Position
		wantErr    bool
	}{
		{
			name:       "Valid input",
			positionID: positionID,
			position: &models.Position{
				ID:        uuid.New(),
				Name:      "C++ Developer",
				Salary:    30999,
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := positionCacheRepo.SetPosition(ctx, tt.positionID, tt.position)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPositionCache_GetPosition(t *testing.T) {
	ctx := context.Background()

	container, connURI := setupRedisContainer(t)
	defer func(container *redisContainer.RedisContainer, ctx context.Context) {
		err := container.Terminate(ctx)
		if err != nil {
			t.Fatal(err)
		}
	}(container, ctx)

	client := redis.NewClient(&redis.Options{
		Addr: connURI,
	})

	positionCacheRepo := NewPositionCache(client)

	positionID := uuid.New()

	err := positionCacheRepo.SetPosition(ctx, positionID.String(), &models.Position{
		ID:        positionID,
		Name:      "Sample",
		Salary:    30999,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})

	require.NoError(t, err)

	tests := []struct {
		name       string
		positionID string
		wantErr    bool
	}{
		{
			name:       "Valid input",
			positionID: positionID.String(),
		},
		{
			name:       "Non-existent position",
			positionID: uuid.New().String(),
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := positionCacheRepo.GetPosition(ctx, tt.positionID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPositionCache_DeletePosition(t *testing.T) {
	ctx := context.Background()

	container, connURI := setupRedisContainer(t)
	defer func(container *redisContainer.RedisContainer, ctx context.Context) {
		err := container.Terminate(ctx)
		if err != nil {
			t.Fatal(err)
		}
	}(container, ctx)

	client := redis.NewClient(&redis.Options{
		Addr: connURI,
	})

	positionCacheRepo := NewPositionCache(client)

	positionID := uuid.New()

	err := positionCacheRepo.SetPosition(ctx, positionID.String(), &models.Position{
		ID:        positionID,
		Name:      "Sample",
		Salary:    30999,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})

	require.NoError(t, err)

	tests := []struct {
		name       string
		positionID string
		wantErr    bool
	}{
		{
			name:       "Valid input",
			positionID: positionID.String(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := positionCacheRepo.DeletePosition(ctx, tt.positionID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
