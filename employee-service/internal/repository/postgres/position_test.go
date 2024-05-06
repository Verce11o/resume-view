package postgres

import (
	"context"
	"errors"
	"github.com/Verce11o/resume-view/employee-service/api"
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
	"testing"
	"time"
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

func setupPostgresContainer(t *testing.T) (*postgres.PostgresContainer, string) {
	ctx := context.Background()
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

func TestPositionRepository_CreatePosition(t *testing.T) {
	ctx := context.Background()

	container, connURI := setupPostgresContainer(t)
	defer container.Terminate(ctx)

	dbPool, err := pgxpool.New(ctx, connURI)
	require.NoError(t, err)

	repo := NewPositionRepository(dbPool)

	positionID := uuid.New()
	tests := []struct {
		name       string
		positionID uuid.UUID
		request    api.CreatePosition
		wantErr    bool
	}{
		{
			name:       "Valid input",
			positionID: positionID,
			request: api.CreatePosition{
				Name:   "Python Developer",
				Salary: 30999,
			},
		},
		{
			name:       "Duplicate position id",
			positionID: positionID,
			request: api.CreatePosition{
				Name:   "Golang Developer",
				Salary: 33333,
			},
			wantErr: true,
		},
		{
			name: "Empty name",
			request: api.CreatePosition{
				Name:   "",
				Salary: 30999,
			},
		},
		{
			name: "Empty salary",
			request: api.CreatePosition{
				Name:   "Python Developer",
				Salary: 0,
			},
		},
		{
			name:    "Empty request",
			request: api.CreatePosition{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := repo.CreatePosition(ctx, tt.positionID, tt.request)
			if tt.wantErr {
				assert.Error(t, err)
			}
		})
	}

}

func TestPositionRepository_GetPosition(t *testing.T) {
	ctx := context.Background()

	container, connURI := setupPostgresContainer(t)
	defer container.Terminate(ctx)

	dbPool, err := pgxpool.New(ctx, connURI)
	require.NoError(t, err)

	repo := NewPositionRepository(dbPool)

	positionID := uuid.New()
	_, err = dbPool.Exec(ctx, "INSERT INTO positions (id, name, salary) VALUES ($1, $2, $3)", positionID, "Sample", 1987)

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
			name:       "Non-existent position",
			positionID: uuid.New(),
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

	container, connURI := setupPostgresContainer(t)
	defer container.Terminate(ctx)

	dbPool, err := pgxpool.New(ctx, connURI)
	require.NoError(t, err)

	repo := NewPositionRepository(dbPool)

	tx, err := dbPool.Begin(ctx)
	require.NoError(t, err)

	for i := 0; i < 10; i++ {
		_, err = tx.Exec(ctx, "INSERT INTO positions (id, name, salary) VALUES ($1, $2, $3)", uuid.New(), "sample", 9999)
		require.NoError(t, err)
	}

	require.NoError(t, tx.Commit(ctx))

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
			name:   "Invalid cursor",
			cursor: "invalid",
			length: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := repo.GetPositionList(ctx, tt.cursor)
			if tt.wantErr {
				assert.Error(t, err)
			}
			assert.Equal(t, len(resp.Positions), tt.length)

		})
	}
}

func TestPositionRepository_UpdatePosition(t *testing.T) {
	ctx := context.Background()

	container, connURI := setupPostgresContainer(t)
	defer container.Terminate(ctx)

	dbPool, err := pgxpool.New(ctx, connURI)
	require.NoError(t, err)

	repo := NewPositionRepository(dbPool)

	positionID := uuid.New()
	_, err = dbPool.Exec(ctx, "INSERT INTO positions (id, name, salary) VALUES ($1, $2, $3)", positionID, "Sample", 1987)

	require.NoError(t, err)

	tests := []struct {
		name       string
		positionID uuid.UUID
		request    api.UpdatePosition
		wantErr    bool
	}{
		{
			name:       "Valid input",
			positionID: positionID,
			request: api.UpdatePosition{
				Name:   "Python Developer",
				Salary: 30999,
			},
		},
		{
			name:       "Non-existing position",
			positionID: uuid.New(),
			request: api.UpdatePosition{
				Name:   "PHP Developer",
				Salary: 9999999,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := repo.UpdatePosition(ctx, tt.positionID, tt.request)
			if tt.wantErr {
				assert.Error(t, err)
			}
		})
	}
}

func TestPositionRepository_DeletePosition(t *testing.T) {
	container, connURI := setupPostgresContainer(t)
	defer container.Terminate(context.Background())

	ctx := context.Background()

	dbPool, err := pgxpool.New(ctx, connURI)
	require.NoError(t, err)

	repo := NewPositionRepository(dbPool)

	positionID := uuid.New()
	_, err = dbPool.Exec(ctx, "INSERT INTO positions (id, name, salary) VALUES ($1, $2, $3)", positionID, "Sample", 1987)

	require.NoError(t, err)

	tests := []struct {
		name       string
		positionID uuid.UUID
		wantErr    bool
	}{
		{
			name:       "Valid position id",
			positionID: positionID,
		},
		{
			name:       "Non-existing position",
			positionID: uuid.New(),
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.DeletePosition(ctx, tt.positionID)
			if tt.wantErr {
				assert.Error(t, err)
			}
		})
	}

}
