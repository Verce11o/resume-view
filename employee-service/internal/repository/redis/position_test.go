package redis

import (
	"context"
	"testing"
	"time"

	"github.com/Verce11o/resume-view/employee-service/internal/lib/customerrors"
	"github.com/Verce11o/resume-view/employee-service/internal/models"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	redisContainer "github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupRedisContainer(ctx context.Context, t *testing.T) (*redisContainer.RedisContainer, string) {
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

func setupPositionRepo(ctx context.Context, t *testing.T) (*PositionCache, *redisContainer.RedisContainer) {
	container, connURI := setupRedisContainer(ctx, t)

	client := redis.NewClient(&redis.Options{
		Addr: connURI,
	})

	positionCacheRepo := NewPositionCache(client)

	return positionCacheRepo, container
}

func TestPositionCache_SetPosition(t *testing.T) {
	ctx := context.Background()

	repo, container := setupPositionRepo(ctx, t)

	defer func(container *redisContainer.RedisContainer, ctx context.Context) {
		err := container.Terminate(ctx)
		if err != nil {
			t.Fatalf("could not terminate redis container: %v", err.Error())
		}
	}(container, ctx)

	positionID := uuid.New().String()

	tests := []struct {
		name       string
		positionID string
		position   *models.Position
		wantErr    error
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
			err := repo.SetPosition(ctx, tt.positionID, tt.position)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestPositionCache_GetPosition(t *testing.T) {
	ctx := context.Background()

	repo, container := setupPositionRepo(ctx, t)

	defer func(container *redisContainer.RedisContainer, ctx context.Context) {
		err := container.Terminate(ctx)
		if err != nil {
			t.Fatalf("could not terminate redis container: %v", err.Error())
		}
	}(container, ctx)

	positionID := uuid.New()

	err := repo.SetPosition(ctx, positionID.String(), &models.Position{
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
		wantErr    error
	}{
		{
			name:       "Valid input",
			positionID: positionID.String(),
		},
		{
			name:       "Non-existent position",
			positionID: uuid.New().String(),
			wantErr:    customerrors.ErrPositionNotCached,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := repo.GetPosition(ctx, tt.positionID)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestPositionCache_DeletePosition(t *testing.T) {
	ctx := context.Background()

	repo, container := setupPositionRepo(ctx, t)

	defer func(container *redisContainer.RedisContainer, ctx context.Context) {
		err := container.Terminate(ctx)
		if err != nil {
			t.Fatalf("could not terminate redis container: %v", err.Error())
		}
	}(container, ctx)

	positionID := uuid.New()

	err := repo.SetPosition(ctx, positionID.String(), &models.Position{
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
		wantErr    error
	}{
		{
			name:       "Valid input",
			positionID: positionID.String(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.DeletePosition(ctx, tt.positionID)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
