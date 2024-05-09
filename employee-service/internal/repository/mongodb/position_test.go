package mongodb

import (
	"context"
	"testing"
	"time"

	"github.com/Verce11o/resume-view/employee-service/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setupMongoDBContainer(t *testing.T) (testcontainers.Container, string) {
	ctx := context.Background()

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

func TestPositionRepository_CreatePosition(t *testing.T) {
	ctx := context.Background()

	container, connURI := setupMongoDBContainer(t)
	defer func(container testcontainers.Container, ctx context.Context) {
		err := container.Terminate(ctx)
		if err != nil {
			t.Fatal(err)
		}
	}(container, ctx)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connURI), options.Client().SetMaxConnIdleTime(3*time.Second))

	require.NoError(t, err)

	repo := NewPositionRepository(client.Database("employees"))

	positionID := uuid.New()

	tests := []struct {
		name    string
		request domain.CreatePosition
		wantErr bool
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
			wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := repo.CreatePosition(ctx, tt.request)
			if tt.wantErr {
				assert.Error(t, err)
			}
		})
	}
}

func TestPositionRepository_GetPosition(t *testing.T) {
	ctx := context.Background()

	container, connURI := setupMongoDBContainer(t)
	defer func(container testcontainers.Container, ctx context.Context) {
		err := container.Terminate(ctx)
		if err != nil {
			t.Fatal(err)
		}
	}(container, ctx)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connURI), options.Client().SetMaxConnIdleTime(3*time.Second))

	require.NoError(t, err)

	err = client.Ping(ctx, nil)

	require.NoError(t, err)

	repo := NewPositionRepository(client.Database("employees"))

	positionID := uuid.New()

	_, err = repo.CreatePosition(ctx, domain.CreatePosition{
		ID:     positionID,
		Name:   "Go Developer",
		Salary: 30999,
	})

	require.NoError(t, err)

	tests := []struct {
		name       string
		positionID uuid.UUID
		wantErr    bool
	}{
		{
			name:       "Valid input",
			positionID: positionID,
		},
		{
			name:       "Non-existent position id",
			positionID: uuid.Nil,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := repo.GetPosition(ctx, tt.positionID)
			if tt.wantErr {
				assert.Error(t, err)
			}
		})
	}
}

func TestPositionRepository_GetPositionList(t *testing.T) {
	ctx := context.Background()

	container, connURI := setupMongoDBContainer(t)
	defer func(container testcontainers.Container, ctx context.Context) {
		err := container.Terminate(ctx)
		if err != nil {
			t.Fatal(err)
		}
	}(container, ctx)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connURI), options.Client().SetMaxConnIdleTime(3*time.Second))

	require.NoError(t, err)

	err = client.Ping(ctx, nil)

	require.NoError(t, err)

	repo := NewPositionRepository(client.Database("employees"))

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
		wantErr bool
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
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := repo.GetPositionList(ctx, tt.cursor)
			if tt.wantErr {
				assert.Error(t, err)
			}
			assert.Equal(t, len(resp.Positions), tt.length)
			nextCursor = resp.Cursor
		})
	}
}

func TestPositionRepository_UpdatePosition(t *testing.T) {
	ctx := context.Background()

	container, connURI := setupMongoDBContainer(t)
	defer func(container testcontainers.Container, ctx context.Context) {
		err := container.Terminate(ctx)
		if err != nil {
			t.Fatal(err)
		}
	}(container, ctx)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connURI), options.Client().SetMaxConnIdleTime(3*time.Second))

	require.NoError(t, err)

	err = client.Ping(ctx, nil)

	require.NoError(t, err)

	repo := NewPositionRepository(client.Database("employees"))

	positionID := uuid.New()

	_, err = repo.CreatePosition(ctx, domain.CreatePosition{
		ID:     positionID,
		Name:   "Sample",
		Salary: 30999,
	})
	require.NoError(t, err)

	tests := []struct {
		name       string
		positionID uuid.UUID
		request    domain.UpdatePosition
		wantErr    bool
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
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err = repo.UpdatePosition(ctx, tt.request)
			if tt.wantErr {
				assert.Error(t, err)
			}
		})
	}
}

func TestPositionRepository_DeletePosition(t *testing.T) {
	ctx := context.Background()

	container, connURI := setupMongoDBContainer(t)
	defer func(container testcontainers.Container, ctx context.Context) {
		err := container.Terminate(ctx)
		if err != nil {
			t.Fatal(err)
		}
	}(container, ctx)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connURI), options.Client().SetMaxConnIdleTime(3*time.Second))

	require.NoError(t, err)

	err = client.Ping(ctx, nil)

	require.NoError(t, err)

	repo := NewPositionRepository(client.Database("employees"))

	positionID := uuid.New()

	_, err = repo.CreatePosition(ctx, domain.CreatePosition{
		ID:     positionID,
		Name:   "Sample",
		Salary: 30999,
	})
	require.NoError(t, err)

	tests := []struct {
		name       string
		positionID uuid.UUID
		wantErr    bool
	}{
		{
			name:       "Valid input",
			positionID: positionID,
		},
		{
			name:       "Non-existent position id",
			positionID: uuid.Nil,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err = repo.DeletePosition(ctx, tt.positionID)
			if tt.wantErr {
				assert.Error(t, err)
			}
		})
	}
}
