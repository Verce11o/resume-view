package postgres

import (
	"context"
	"testing"

	"github.com/Verce11o/resume-view/employee-service/api"
	_ "github.com/flashlabs/rootpath"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestEmployeeRepository_CreateEmployee(t *testing.T) {
	ctx := context.Background()

	container, connURI := setupPostgresContainer(t)
	defer func(container *postgres.PostgresContainer, ctx context.Context) {
		err := container.Terminate(ctx)
		if err != nil {
			t.Fatal(err)
		}
	}(container, ctx)

	dbPool, err := pgxpool.New(ctx, connURI)
	require.NoError(t, err)

	repo := NewEmployeeRepository(dbPool)

	employeeID := uuid.New()
	positionID := uuid.New()

	tests := []struct {
		name       string
		employeeID uuid.UUID
		positionID uuid.UUID
		request    api.CreateEmployee
		wantErr    bool
	}{
		{
			name:       "Valid input",
			employeeID: employeeID,
			positionID: positionID,
			request: api.CreateEmployee{
				FirstName:    "John",
				LastName:     "Doe",
				PositionName: "Go Developer",
				Salary:       12345,
			},
		},
		{
			name:       "Invalid position id",
			positionID: uuid.Nil,
			employeeID: employeeID,
			request: api.CreateEmployee{
				FirstName:    "John",
				LastName:     "Doe",
				PositionName: "Python Developer",
				Salary:       12345,
			},
			wantErr: true,
		},
		{
			name:       "Invalid employee id",
			employeeID: uuid.Nil,
			positionID: positionID,
			request: api.CreateEmployee{
				FirstName:    "John",
				LastName:     "Doe",
				PositionName: "Python Developer",
				Salary:       12345,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := repo.CreateEmployee(ctx, tt.employeeID, tt.positionID, tt.request)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestEmployeeRepository_GetEmployee(t *testing.T) {
	ctx := context.Background()

	container, connURI := setupPostgresContainer(t)
	defer func(container *postgres.PostgresContainer, ctx context.Context) {
		err := container.Terminate(ctx)
		if err != nil {
			t.Fatal(err)
		}
	}(container, ctx)

	dbPool, err := pgxpool.New(ctx, connURI)
	require.NoError(t, err)

	repo := NewEmployeeRepository(dbPool)

	employeeID := uuid.New()
	positionID := uuid.New()

	_, err = dbPool.Exec(ctx, "INSERT INTO positions (id, name, salary)  VALUES  ($1, $2, $3)", positionID, "Go Developer", 123456)

	require.NoError(t, err)

	_, err = dbPool.Exec(ctx, "INSERT INTO employees (id, first_name, last_name, position_id) VALUES ($1, $2, $3, $4)", employeeID, "John", "Doe", positionID)

	require.NoError(t, err)

	tests := []struct {
		name       string
		employeeID uuid.UUID
		wantErr    bool
	}{
		{
			name:       "Valid input",
			employeeID: employeeID,
		},
		{
			name:       "Non-existent employee",
			employeeID: uuid.New(),
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := repo.GetEmployee(ctx, tt.employeeID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestEmployeeRepository_GetEmployeeList(t *testing.T) {
	ctx := context.Background()

	container, connURI := setupPostgresContainer(t)
	defer func(container *postgres.PostgresContainer, ctx context.Context) {
		err := container.Terminate(ctx)
		if err != nil {
			t.Fatal(err)
		}
	}(container, ctx)

	dbPool, err := pgxpool.New(ctx, connURI)
	require.NoError(t, err)

	repo := NewEmployeeRepository(dbPool)

	tx, err := dbPool.Begin(ctx)
	require.NoError(t, err)

	positionID := uuid.New()
	_, err = dbPool.Exec(ctx, "INSERT INTO positions (id, name, salary) VALUES ($1, $2, $3)", positionID, "Sample", 1987)

	require.NoError(t, err)

	for i := 0; i < 10; i++ {
		_, err = tx.Exec(ctx, "INSERT INTO employees (id, first_name, last_name, position_id) VALUES ($1, $2, $3, $4)", uuid.New(), "John", "Doe", positionID)
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
			name:    "Invalid cursor",
			cursor:  "invalid",
			length:  0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := repo.GetEmployeeList(ctx, tt.cursor)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, len(resp.Employees), tt.length)
			nextCursor = resp.Cursor
		})
	}
}

func TestEmployeeRepository_UpdateEmployee(t *testing.T) {
	ctx := context.Background()

	container, connURI := setupPostgresContainer(t)
	defer func(container *postgres.PostgresContainer, ctx context.Context) {
		err := container.Terminate(ctx)
		if err != nil {
			t.Fatal(err)
		}
	}(container, ctx)

	dbPool, err := pgxpool.New(ctx, connURI)
	require.NoError(t, err)

	repo := NewEmployeeRepository(dbPool)

	positionID := uuid.New()
	_, err = dbPool.Exec(ctx, "INSERT INTO positions (id, name, salary) VALUES ($1, $2, $3)", positionID, "Sample", 1987)

	require.NoError(t, err)

	employeeID := uuid.New()
	_, err = dbPool.Exec(ctx, "INSERT INTO employees (id, first_name, last_name, position_id) VALUES ($1, $2, $3, $4)", employeeID, "John", "Doe", positionID)

	require.NoError(t, err)

	tests := []struct {
		name       string
		employeeID uuid.UUID
		positionID uuid.UUID
		request    api.UpdateEmployee
		wantErr    bool
	}{
		{
			name:       "Valid input",
			employeeID: employeeID,
			positionID: positionID,
			request: api.UpdateEmployee{
				FirstName:  "NewName",
				LastName:   "NewLastName",
				PositionId: positionID.String(),
			},
		},
		{
			name:       "Non-existing position",
			employeeID: employeeID,
			request: api.UpdateEmployee{
				FirstName:  "NewName",
				LastName:   "NewLastName",
				PositionId: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := repo.UpdateEmployee(ctx, tt.employeeID, tt.request)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestEmployeeRepository_DeleteEmployee(t *testing.T) {
	ctx := context.Background()

	container, connURI := setupPostgresContainer(t)
	defer func(container *postgres.PostgresContainer, ctx context.Context) {
		err := container.Terminate(ctx)
		if err != nil {
			t.Fatal(err)
		}
	}(container, ctx)

	dbPool, err := pgxpool.New(ctx, connURI)
	require.NoError(t, err)

	repo := NewEmployeeRepository(dbPool)

	positionID := uuid.New()
	_, err = dbPool.Exec(ctx, "INSERT INTO positions (id, name, salary) VALUES ($1, $2, $3)", positionID, "Sample", 1987)

	require.NoError(t, err)

	employeeID := uuid.New()
	_, err = dbPool.Exec(ctx, "INSERT INTO employees (id, first_name, last_name, position_id) VALUES ($1, $2, $3, $4)", employeeID, "John", "Doe", positionID)

	require.NoError(t, err)

	tests := []struct {
		name       string
		employeeID uuid.UUID
		wantErr    bool
	}{
		{
			name:       "Valid employee id",
			employeeID: employeeID,
		},
		{
			name:       "Non-existing employee",
			employeeID: uuid.New(),
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.DeleteEmployee(ctx, tt.employeeID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
