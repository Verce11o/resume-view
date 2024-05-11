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
	redisContainer "github.com/testcontainers/testcontainers-go/modules/redis"
)

func setupEmployeeRepo(ctx context.Context, t *testing.T) (*EmployeeCache, *redisContainer.RedisContainer) {
	container, connURI := setupRedisContainer(ctx, t)

	client := redis.NewClient(&redis.Options{
		Addr: connURI,
	})

	employeeCacheRepo := NewEmployeeCache(client)

	return employeeCacheRepo, container

}

func TestEmployeeCache_SetEmployee(t *testing.T) {
	ctx := context.Background()

	repo, container := setupEmployeeRepo(ctx, t)

	defer func(container *redisContainer.RedisContainer, ctx context.Context) {
		err := container.Terminate(ctx)
		if err != nil {
			t.Fatalf("could not terminate redis container: %v", err.Error())
		}
	}(container, ctx)

	employeeID := uuid.New().String()

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
			err := repo.SetEmployee(ctx, tt.employeeID, tt.employee)
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

	repo, container := setupEmployeeRepo(ctx, t)

	defer func(container *redisContainer.RedisContainer, ctx context.Context) {
		err := container.Terminate(ctx)
		if err != nil {
			t.Fatalf("could not terminate redis container: %v", err.Error())
		}
	}(container, ctx)

	employeeID := uuid.New()

	err := repo.SetEmployee(ctx, employeeID.String(), &models.Employee{
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
		wantErr    error
	}{
		{
			name:       "Valid input",
			employeeID: employeeID.String(),
		},
		{
			name:       "Non-existent employee",
			employeeID: uuid.New().String(),
			wantErr:    customerrors.ErrEmployeeNotCached,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := repo.GetEmployee(ctx, tt.employeeID)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestEmployeeCache_DeleteEmployee(t *testing.T) {
	ctx := context.Background()

	repo, container := setupEmployeeRepo(ctx, t)

	defer func(container *redisContainer.RedisContainer, ctx context.Context) {
		err := container.Terminate(ctx)
		if err != nil {
			t.Fatalf("could not terminate redis container: %v", err.Error())
		}
	}(container, ctx)

	employeeID := uuid.New()

	err := repo.SetEmployee(ctx, employeeID.String(), &models.Employee{
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
		wantErr    error
	}{
		{
			name:       "Valid input",
			employeeID: employeeID.String(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.DeleteEmployee(ctx, tt.employeeID)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
