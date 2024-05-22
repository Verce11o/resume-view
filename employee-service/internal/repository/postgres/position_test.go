//go:build integration

package postgres

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Verce11o/resume-view/employee-service/internal/domain"
	"github.com/Verce11o/resume-view/employee-service/internal/lib/customerrors"
	"github.com/Verce11o/resume-view/employee-service/internal/models"
	_ "github.com/flashlabs/rootpath"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func runMigrations(t *testing.T, connURI string) {
	m, err := migrate.New(
		"file://migrations",
		connURI)

	require.NoError(t, err)

	defer m.Close()

	err = m.Up()

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		require.NoError(t, err)
	}
}

func setupPostgresContainer(ctx context.Context, t *testing.T) (*postgres.PostgresContainer, string) {
	postgresContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:latest"),
		postgres.WithDatabase("employees"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("vercello"),
		testcontainers.WithWaitStrategy(
			wait.
				ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(3*time.Second),
		),
	)
	require.NoError(t, err)

	connURI, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	runMigrations(t, connURI)

	return postgresContainer, connURI
}

type PositionRepositorySuite struct {
	suite.Suite
	ctx       context.Context
	repo      *PositionRepository
	container *postgres.PostgresContainer
}

func (p *PositionRepositorySuite) SetupSuite() {
	p.ctx = context.Background()

	container, connURI := setupPostgresContainer(p.ctx, p.T())
	dbPool, err := pgxpool.New(p.ctx, connURI)
	require.NoError(p.T(), err)

	positionRepo := NewPositionRepository(dbPool)

	p.repo = positionRepo
	p.container = container
}

func (p *PositionRepositorySuite) TearDownSuite() {
	err := p.container.Terminate(p.ctx)
	if err != nil {
		p.T().Fatalf("could not terminate postgres container: %v", err.Error())
	}
}

func (p *PositionRepositorySuite) TestCreatePosition() {
	positionID := uuid.New()
	tests := []struct {
		name     string
		request  domain.CreatePosition
		response models.Position
		wantErr  error
	}{
		{
			name: "Valid input",
			request: domain.CreatePosition{
				ID:     positionID,
				Name:   "Python Developer",
				Salary: 30999,
			},
			response: models.Position{
				ID:     positionID,
				Name:   "Python Developer",
				Salary: 30999,
			},
		},
		{
			name: "Duplicate position id",
			request: domain.CreatePosition{
				ID:     positionID,
				Name:   "Golang Developer",
				Salary: 33333,
			},
			response: models.Position{},
			wantErr:  customerrors.ErrDuplicateID,
		},
	}

	for _, tt := range tests {
		p.Run(tt.name, func() {
			resp, err := p.repo.CreatePosition(p.ctx, tt.request)
			assert.ErrorIs(p.T(), err, tt.wantErr)
			assert.EqualExportedValues(p.T(), tt.response, resp)
		})
	}
}

func (p *PositionRepositorySuite) TestGetPosition() {
	positionID := uuid.New()

	position, err := p.repo.CreatePosition(p.ctx, domain.CreatePosition{
		ID:     positionID,
		Name:   "Go Developer",
		Salary: 30999,
	})

	require.NoError(p.T(), err)

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
			name:       "Non-existent position",
			positionID: uuid.New(),
			wantErr:    customerrors.ErrPositionNotFound,
		},
	}

	for _, tt := range tests {
		p.Run(tt.name, func() {
			resp, err := p.repo.GetPosition(p.ctx, tt.positionID)
			if tt.wantErr != nil {
				assert.ErrorIs(p.T(), err, tt.wantErr)
				assert.NotEqual(p.T(), resp, position)

				return
			}

			assert.Equal(p.T(), resp, position)
		})
	}
}

func (p *PositionRepositorySuite) TestGetPositionList() {
	for i := 0; i < 10; i++ {
		_, err := p.repo.CreatePosition(p.ctx, domain.CreatePosition{
			ID:     uuid.New(),
			Name:   "Sample",
			Salary: 30999,
		})
		require.NoError(p.T(), err)
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
		p.Run(tt.name, func() {
			resp, err := p.repo.GetPositionList(p.ctx, tt.cursor)
			assert.ErrorIs(p.T(), err, tt.wantErr)

			assert.Equal(p.T(), len(resp.Positions), tt.length)
			nextCursor = resp.Cursor
		})
	}
}

func (p *PositionRepositorySuite) TestUpdatePosition() {
	positionID := uuid.New()

	_, err := p.repo.CreatePosition(p.ctx, domain.CreatePosition{
		ID:     positionID,
		Name:   "Go developer",
		Salary: 30999,
	})

	require.NoError(p.T(), err)

	tests := []struct {
		name       string
		positionID uuid.UUID
		request    domain.UpdatePosition
		response   models.Position
		wantErr    error
	}{
		{
			name: "Valid input",
			request: domain.UpdatePosition{
				ID:     positionID,
				Name:   "Python Developer",
				Salary: 30999,
			},
			response: models.Position{
				ID:     positionID,
				Name:   "Python Developer",
				Salary: 30999,
			},
		},
		{
			name: "Non-existing position",
			request: domain.UpdatePosition{
				ID:     uuid.New(),
				Name:   "PHP Developer",
				Salary: 9999999,
			},
			wantErr: customerrors.ErrPositionNotFound,
		},
	}

	for _, tt := range tests {
		p.Run(tt.name, func() {
			resp, err := p.repo.UpdatePosition(p.ctx, tt.request)
			assert.ErrorIs(p.T(), err, tt.wantErr)
			assert.EqualExportedValues(p.T(), tt.response, resp)
		})
	}
}

func (p *PositionRepositorySuite) TestDeletePosition() {
	positionID := uuid.New()

	_, err := p.repo.CreatePosition(p.ctx, domain.CreatePosition{
		ID:     positionID,
		Name:   "Go developer",
		Salary: 30999,
	})

	require.NoError(p.T(), err)

	tests := []struct {
		name       string
		positionID uuid.UUID
		wantErr    error
	}{
		{
			name:       "Valid position id",
			positionID: positionID,
		},
		{
			name:       "Non-existing position",
			positionID: uuid.New(),
			wantErr:    customerrors.ErrPositionNotFound,
		},
	}
	for _, tt := range tests {
		p.Run(tt.name, func() {
			err := p.repo.DeletePosition(p.ctx, tt.positionID)
			assert.ErrorIs(p.T(), err, tt.wantErr)

			if tt.wantErr != nil {
				_, err = p.repo.GetPosition(p.ctx, tt.positionID)
				assert.ErrorIs(p.T(), err, customerrors.ErrPositionNotFound)
			}
		})
	}
}

func TestPositionRepositorySuite(t *testing.T) {
	suite.Run(t, new(PositionRepositorySuite))
}
