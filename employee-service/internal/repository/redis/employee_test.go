package redis

import (
	"context"
	"github.com/Verce11o/resume-view/employee-service/internal/models"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	redisContainer "github.com/testcontainers/testcontainers-go/modules/redis"
	"testing"
	"time"
)

func TestEmployeeCache_SetEmployee(t *testing.T) {

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

	employeeID := uuid.New().String()
	employeeCacheRepo := NewEmployeeCache(client)

	tests := []struct {
		name       string
		employeeID string
		employee   *models.Employee
		wantErr    bool
	}{
		{
			name:       "Valid input",
			employeeID: employeeID,
			employee: &models.Employee{
				ID:         uuid.New(),
				FirstName:  "John",
				LastName:   "Doe",
				PositionID: uuid.New(),
				CreatedAt:  time.Now().UTC(),
				UpdatedAt:  time.Now().UTC(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := employeeCacheRepo.SetEmployee(ctx, tt.employeeID, tt.employee)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestEmployeeCache_GetEmployee(t *testing.T) {
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

	employeeCacheRepo := NewEmployeeCache(client)

	employeeID := uuid.New()

	err := employeeCacheRepo.SetEmployee(ctx, employeeID.String(), &models.Employee{
		ID:        employeeID,
		FirstName: "John",
		LastName:  "Doe",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})

	require.NoError(t, err)

	tests := []struct {
		name       string
		employeeID string
		wantErr    bool
	}{
		{
			name:       "Valid input",
			employeeID: employeeID.String(),
		},
		{
			name:       "Non-existent employee",
			employeeID: uuid.New().String(),
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := employeeCacheRepo.GetEmployee(ctx, tt.employeeID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestEmployeeCache_DeleteEmployee(t *testing.T) {
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

	employeeCacheRepo := NewEmployeeCache(client)

	employeeID := uuid.New()

	err := employeeCacheRepo.SetEmployee(ctx, employeeID.String(), &models.Employee{
		ID:         employeeID,
		FirstName:  "John",
		LastName:   "Doe",
		PositionID: uuid.New(),
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	})

	require.NoError(t, err)

	tests := []struct {
		name       string
		employeeID string
		wantErr    bool
	}{
		{
			name:       "Valid input",
			employeeID: employeeID.String(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := employeeCacheRepo.DeleteEmployee(ctx, tt.employeeID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
