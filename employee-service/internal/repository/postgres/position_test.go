package postgres

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Verce11o/resume-view/employee-service/internal/domain"
	"github.com/Verce11o/resume-view/employee-service/internal/lib/customerrors"
	_ "github.com/flashlabs/rootpath"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func setupPositionRepo(ctx context.Context, t *testing.T) (*PositionRepository, *postgres.PostgresContainer) {
	container, connURI := setupPostgresContainer(ctx, t)

	dbPool, err := pgxpool.New(ctx, connURI)
	require.NoError(t, err)

	repo := NewPositionRepository(dbPool)

	return repo, container
}

func TestPositionRepository_CreatePosition(t *testing.T) {
	ctx := context.Background()

	repo, container := setupPositionRepo(ctx, t)
	defer func(container *postgres.PostgresContainer, ctx context.Context) {
		err := container.Terminate(ctx)
		if err != nil {
			t.Fatalf("could not terminate postgres container: %v", err.Error())
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
			wantErr: customerrors.ErrDuplicateID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := repo.CreatePosition(ctx, tt.request)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestPositionRepository_GetPosition(t *testing.T) {
	ctx := context.Background()

	repo, container := setupPositionRepo(ctx, t)
	defer func(container *postgres.PostgresContainer, ctx context.Context) {
		err := container.Terminate(ctx)
		if err != nil {
			t.Fatalf("could not terminate postgres container: %v", err.Error())
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
			name:       "Non-existent position",
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
	defer func(container *postgres.PostgresContainer, ctx context.Context) {
		err := container.Terminate(ctx)
		if err != nil {
			t.Fatalf("could not terminate postgres container: %v", err.Error())
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
	defer func(container *postgres.PostgresContainer, ctx context.Context) {
		err := container.Terminate(ctx)
		if err != nil {
			t.Fatalf("could not terminate postgres container: %v", err.Error())
		}
	}(container, ctx)

	positionID := uuid.New()

	_, err := repo.CreatePosition(ctx, domain.CreatePosition{
		ID:     positionID,
		Name:   "Go developer",
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
		t.Run(tt.name, func(t *testing.T) {
			_, err := repo.UpdatePosition(ctx, tt.request)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestPositionRepository_DeletePosition(t *testing.T) {
	ctx := context.Background()

	repo, container := setupPositionRepo(ctx, t)
	defer func(container *postgres.PostgresContainer, ctx context.Context) {
		err := container.Terminate(ctx)
		if err != nil {
			t.Fatalf("could not terminate postgres container: %v", err.Error())
		}
	}(container, ctx)

	positionID := uuid.New()

	_, err := repo.CreatePosition(ctx, domain.CreatePosition{
		ID:     positionID,
		Name:   "Go developer",
		Salary: 30999,
	})

	require.NoError(t, err)

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
		t.Run(tt.name, func(t *testing.T) {
			err := repo.DeletePosition(ctx, tt.positionID)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
