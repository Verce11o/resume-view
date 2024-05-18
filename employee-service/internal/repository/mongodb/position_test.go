//go:build integration

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
	"github.com/stretchr/testify/suite"
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

type PositionRepositorySuite struct {
	suite.Suite
	ctx       context.Context
	client    *mongo.Client
	container testcontainers.Container
	repo      *PositionRepository
}

func (s *PositionRepositorySuite) SetupSuite() {
	s.ctx = context.Background()
	container, connURI := setupMongoDBContainer(s.ctx, s.T())

	client, err := mongo.Connect(s.ctx,
		options.Client().ApplyURI(connURI),
		options.Client().SetMaxConnIdleTime(3*time.Second))
	require.NoError(s.T(), err)

	s.repo = NewPositionRepository(client.Database("employees"))
	s.client = client
	s.container = container
}

func (s *PositionRepositorySuite) TearDownSuite() {
	err := s.container.Terminate(s.ctx)
	require.NoError(s.T(), err)

}

func (s *PositionRepositorySuite) TestCreatePosition() {
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
		s.Run(tt.name, func() {

			_, err := s.repo.CreatePosition(s.ctx, tt.request)
			assert.ErrorIs(s.T(), err, tt.wantErr)
		})
	}
}

func (s *PositionRepositorySuite) TestGetPosition() {

	positionID := uuid.New()

	position, err := s.repo.CreatePosition(s.ctx, domain.CreatePosition{
		ID:     positionID,
		Name:   "Go Developer",
		Salary: 30999,
	})
	require.NoError(s.T(), err)

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
		s.Run(tt.name, func() {

			resp, err := s.repo.GetPosition(s.ctx, tt.positionID)
			if tt.wantErr != nil {
				assert.ErrorIs(s.T(), err, tt.wantErr)
				assert.NotEqual(s.T(), resp, position)
				return
			}
			assert.Equal(s.T(), resp, position)
		})
	}
}

func (s *PositionRepositorySuite) TestGetPositionList() {

	for i := 0; i < 10; i++ {
		_, err := s.repo.CreatePosition(s.ctx, domain.CreatePosition{
			ID:     uuid.New(),
			Name:   "Sample",
			Salary: 30999,
		})
		require.NoError(s.T(), err)
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
		s.Run(tt.name, func() {

			resp, err := s.repo.GetPositionList(s.ctx, tt.cursor)
			assert.ErrorIs(s.T(), err, tt.wantErr)
			assert.Equal(s.T(), len(resp.Positions), tt.length)
			nextCursor = resp.Cursor
		})
	}
}

func (s *PositionRepositorySuite) TestUpdatePosition() {

	positionID := uuid.New()

	_, err := s.repo.CreatePosition(s.ctx, domain.CreatePosition{
		ID:     positionID,
		Name:   "Sample",
		Salary: 30999,
	})
	require.NoError(s.T(), err)

	tests := []struct {
		name    string
		request domain.UpdatePosition
		wantErr error
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
		s.Run(tt.name, func() {

			_, err := s.repo.UpdatePosition(s.ctx, tt.request)
			assert.ErrorIs(s.T(), err, tt.wantErr)
		})
	}
}

func (s *PositionRepositorySuite) TestDeletePosition() {

	positionID := uuid.New()

	_, err := s.repo.CreatePosition(s.ctx, domain.CreatePosition{
		ID:     positionID,
		Name:   "Sample",
		Salary: 30999,
	})
	require.NoError(s.T(), err)

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
		s.Run(tt.name, func() {

			err := s.repo.DeletePosition(s.ctx, tt.positionID)
			assert.ErrorIs(s.T(), err, tt.wantErr)
		})
	}
}

func TestPositionRepositorySuite(t *testing.T) {
	suite.Run(t, new(PositionRepositorySuite))
}
