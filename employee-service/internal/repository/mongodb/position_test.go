package mongodb

import (
	"context"
	"testing"
	"time"

	"github.com/Verce11o/resume-view/employee-service/internal/domain"
	"github.com/Verce11o/resume-view/employee-service/internal/lib/customerrors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setupMongoDBContainer(ctx context.Context, t *testing.T) (testcontainers.Container, string) {
	entrypointScript := []string{
		"/bin/bash", "-c",
		`echo "rs.initiate()" > /docker-entrypoint-initdb.d/1-init-replicaset.js &&
		exec /usr/local/bin/docker-entrypoint.sh mongod --replSet rs0 --bind_ip_all --noauth`,
	}

	mongoContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "mongo:latest",
			ExposedPorts: []string{"27017/tcp"},
			Env: map[string]string{
				"MONGO_APP_DATABASE": "employees",
				"MONGO_REPLICA_PORT": "27018",
			},
			WaitingFor: wait.ForListeningPort("27017/tcp"),
			Entrypoint: entrypointScript,
		},
		Started: true,
	})
	require.NoError(t, err)

	connURI, err := mongoContainer.Endpoint(ctx, "mongodb")
	require.NoError(t, err)

	return mongoContainer, connURI + "/?directConnection=true&tls=false"
}

func setupPositionRepo(ctx context.Context, t *testing.T) (*PositionRepository, testcontainers.Container) {
	container, connURI := setupMongoDBContainer(ctx, t)

	client, err := mongo.Connect(ctx,
		options.Client().ApplyURI(connURI),
		options.Client().SetMaxConnIdleTime(3*time.Second))

	require.NoError(t, err)

	repo := NewPositionRepository(client.Database("employees"))

	return repo, container
}

func TestPositionRepository_CreatePosition(t *testing.T) {
	ctx := context.Background()

	repo, container := setupPositionRepo(ctx, t)
	defer func(container testcontainers.Container, ctx context.Context) {
		err := container.Terminate(ctx)
		if err != nil {
			t.Fatalf("could not terminate mongo container: %v", err.Error())
		}
	}(container, ctx)

	positionID := uuid.New()

	tests := []struct {
		name    string
		request domain.CreatePosition
		wantErr error
	}{
		{
			name: "Valid input",
			request: domain.CreatePosition{
				ID:     positionID,
				Name:   "Go Developer",
				Salary: 30999,
			},
		},
		{
			name: "Duplicate position id",
			request: domain.CreatePosition{
				ID:     positionID,
				Name:   "Go Developer",
				Salary: 30999,
			},
			wantErr: customerrors.ErrDuplicateID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := repo.CreatePosition(ctx, tt.request)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestPositionRepository_GetPosition(t *testing.T) {
	ctx := context.Background()

	repo, container := setupPositionRepo(ctx, t)
	defer func(container testcontainers.Container, ctx context.Context) {
		err := container.Terminate(ctx)
		if err != nil {
			t.Fatalf("could not terminate mongo container: %v", err.Error())
		}
	}(container, ctx)

	positionID := uuid.New()

	position, err := repo.CreatePosition(ctx, domain.CreatePosition{
		ID:     positionID,
		Name:   "Go Developer",
		Salary: 30999,
	})

	require.NoError(t, err)

	tests := []struct {
		name       string
		positionID uuid.UUID
		wantErr    error
	}{
		{
			name:       "Valid input",
			positionID: positionID,
		},
		{
			name:       "Non-existent position id",
			positionID: uuid.New(),
			wantErr:    customerrors.ErrPositionNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := repo.GetPosition(ctx, tt.positionID)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.NotEqual(t, resp, position)

				return
			}

			assert.Equal(t, resp, position)
		})
	}
}

func TestPositionRepository_GetPositionList(t *testing.T) {
	ctx := context.Background()

	repo, container := setupPositionRepo(ctx, t)
	defer func(container testcontainers.Container, ctx context.Context) {
		err := container.Terminate(ctx)
		if err != nil {
			t.Fatalf("could not terminate mongo container: %v", err.Error())
		}
	}(container, ctx)

	for i := 0; i < 10; i++ {
		_, err := repo.CreatePosition(ctx, domain.CreatePosition{
			ID:     uuid.New(),
			Name:   "Sample",
			Salary: 30999,
		})
		require.NoError(t, err)
	}

	var nextCursor string
	tests := []struct {
		name    string
		cursor  string
		length  int
		wantErr error
	}{
		{
			name:   "First page",
			cursor: nextCursor,
			length: 5,
		},
		{
			name:   "Second page",
			cursor: nextCursor,
			length: 5,
		},
		{
			name:    "Invalid cursor",
			cursor:  "invalid",
			length:  0,
			wantErr: customerrors.ErrInvalidCursor,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := repo.GetPositionList(ctx, tt.cursor)
			assert.ErrorIs(t, err, tt.wantErr)

			assert.Equal(t, len(resp.Positions), tt.length)

			nextCursor = resp.Cursor
		})
	}
}

func TestPositionRepository_UpdatePosition(t *testing.T) {
	ctx := context.Background()

	repo, container := setupPositionRepo(ctx, t)
	defer func(container testcontainers.Container, ctx context.Context) {
		err := container.Terminate(ctx)
		if err != nil {
			t.Fatalf("could not terminate mongo container: %v", err.Error())
		}
	}(container, ctx)

	positionID := uuid.New()

	_, err := repo.CreatePosition(ctx, domain.CreatePosition{
		ID:     positionID,
		Name:   "Sample",
		Salary: 30999,
	})
	require.NoError(t, err)

	tests := []struct {
		name       string
		positionID uuid.UUID
		request    domain.UpdatePosition
		wantErr    error
	}{
		{
			name: "Valid input",
			request: domain.UpdatePosition{
				ID:     positionID,
				Name:   "NewName",
				Salary: 10300,
			},
		},
		{
			name: "Non-existent position id",
			request: domain.UpdatePosition{
				ID:     uuid.Nil,
				Name:   "NewName",
				Salary: 10300,
			},
			wantErr: customerrors.ErrPositionNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err = repo.UpdatePosition(ctx, tt.request)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestPositionRepository_DeletePosition(t *testing.T) {
	ctx := context.Background()

	repo, container := setupPositionRepo(ctx, t)
	defer func(container testcontainers.Container, ctx context.Context) {
		err := container.Terminate(ctx)
		if err != nil {
			t.Fatalf("could not terminate mongo container: %v", err.Error())
		}
	}(container, ctx)

	positionID := uuid.New()

	_, err := repo.CreatePosition(ctx, domain.CreatePosition{
		ID:     positionID,
		Name:   "Sample",
		Salary: 30999,
	})
	require.NoError(t, err)

	tests := []struct {
		name       string
		positionID uuid.UUID
		wantErr    error
	}{
		{
			name:       "Valid input",
			positionID: positionID,
		},
		{
			name:       "Non-existent position id",
			positionID: uuid.Nil,
			wantErr:    customerrors.ErrPositionNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err = repo.DeletePosition(ctx, tt.positionID)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
