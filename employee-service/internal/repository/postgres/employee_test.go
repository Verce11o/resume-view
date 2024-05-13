package postgres

import (
	"context"
	"testing"

	"github.com/Verce11o/resume-view/employee-service/internal/domain"
	"github.com/Verce11o/resume-view/employee-service/internal/lib/customerrors"
	_ "github.com/flashlabs/rootpath"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func setupEmployeeRepo(ctx context.Context, t *testing.T) (*EmployeeRepository, *postgres.PostgresContainer) {
	container, connURI := setupPostgresContainer(ctx, t)

	dbPool, err := pgxpool.New(ctx, connURI)
	require.NoError(t, err)

	repo := NewEmployeeRepository(dbPool)

	return repo, container
}

func TestEmployeeRepository_CreateEmployee(t *testing.T) {
	ctx := context.Background()

	repo, container := setupEmployeeRepo(ctx, t)
	defer func(container *postgres.PostgresContainer, ctx context.Context) {
		err := container.Terminate(ctx)
		if err != nil {
			t.Fatalf("could not terminate postgres container: %v", err.Error())
		}
	}(container, ctx)

	employeeID := uuid.New()
	positionID := uuid.New()

	tests := []struct {
		name    string
		request domain.CreateEmployee
		wantErr error
	}{
		{
			name: "Valid input",
			request: domain.CreateEmployee{
				EmployeeID:   employeeID,
				PositionID:   positionID,
				FirstName:    "John",
				LastName:     "Doe",
				PositionName: "Go Developer",
				Salary:       12345,
			},
		},
		{
			name: "Duplicate position id",
			request: domain.CreateEmployee{
				EmployeeID:   uuid.New(),
				PositionID:   positionID,
				FirstName:    "John",
				LastName:     "Doe",
				PositionName: "Python Developer",
				Salary:       12345,
			},
			wantErr: customerrors.ErrDuplicateID,
		},
		{
			name: "Duplicate employee id",
			request: domain.CreateEmployee{
				EmployeeID:   employeeID,
				PositionID:   uuid.New(),
				FirstName:    "John",
				LastName:     "Doe",
				PositionName: "Python Developer",
				Salary:       12345,
			},
			wantErr: customerrors.ErrDuplicateID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := repo.CreateEmployee(ctx, tt.request)

			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestEmployeeRepository_GetEmployee(t *testing.T) {
	ctx := context.Background()

	repo, container := setupEmployeeRepo(ctx, t)
	defer func(container *postgres.PostgresContainer, ctx context.Context) {
		err := container.Terminate(ctx)
		if err != nil {
			t.Fatalf("could not terminate postgres container: %v", err.Error())
		}
	}(container, ctx)

	employeeID := uuid.New()
	positionID := uuid.New()

	employee, err := repo.CreateEmployee(ctx, domain.CreateEmployee{
		EmployeeID:   employeeID,
		PositionID:   positionID,
		FirstName:    "John",
		LastName:     "Doe",
		PositionName: "Go Developer",
		Salary:       0,
	})

	require.NoError(t, err)

	tests := []struct {
		name       string
		employeeID uuid.UUID
		wantErr    error
	}{
		{
			name:       "Valid input",
			employeeID: employeeID,
		},
		{
			name:       "Non-existent employee",
			employeeID: uuid.New(),
			wantErr:    customerrors.ErrEmployeeNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := repo.GetEmployee(ctx, tt.employeeID)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.NotEqual(t, resp, employee)

				return
			}

			assert.Equal(t, resp, employee)
		})
	}
}

func TestEmployeeRepository_GetEmployeeList(t *testing.T) {
	ctx := context.Background()

	repo, container := setupEmployeeRepo(ctx, t)
	defer func(container *postgres.PostgresContainer, ctx context.Context) {
		err := container.Terminate(ctx)
		if err != nil {
			t.Fatalf("could not terminate postgres container: %v", err.Error())
		}
	}(container, ctx)

	for i := 0; i < 10; i++ {
		_, err := repo.CreateEmployee(ctx, domain.CreateEmployee{
			EmployeeID:   uuid.New(),
			PositionID:   uuid.New(),
			FirstName:    "John",
			LastName:     "Doe",
			PositionName: "Go Developer",
			Salary:       30999,
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
			resp, err := repo.GetEmployeeList(ctx, tt.cursor)
			assert.ErrorIs(t, err, tt.wantErr)

			assert.Equal(t, len(resp.Employees), tt.length)
			nextCursor = resp.Cursor
		})
	}
}

func TestEmployeeRepository_UpdateEmployee(t *testing.T) {
	ctx := context.Background()
	repo, container := setupEmployeeRepo(ctx, t)
	defer func(container *postgres.PostgresContainer, ctx context.Context) {
		err := container.Terminate(ctx)
		if err != nil {
			t.Fatalf("could not terminate postgres container: %v", err.Error())
		}
	}(container, ctx)

	employeeID := uuid.New()
	positionID := uuid.New()

	_, err := repo.CreateEmployee(ctx, domain.CreateEmployee{
		EmployeeID:   employeeID,
		PositionID:   positionID,
		FirstName:    "John",
		LastName:     "Doe",
		PositionName: "Go developer",
		Salary:       30999,
	})

	require.NoError(t, err)

	tests := []struct {
		name    string
		request domain.UpdateEmployee
		wantErr error
	}{
		{
			name: "Valid input",
			request: domain.UpdateEmployee{
				EmployeeID: employeeID,
				PositionID: positionID,
				FirstName:  "NewName",
				LastName:   "NewLastName",
			},
		},
		{
			name: "Non-existing position",
			request: domain.UpdateEmployee{
				EmployeeID: employeeID,
				PositionID: uuid.Nil,
				FirstName:  "NewName",
				LastName:   "NewLastName",
			},
			wantErr: customerrors.ErrPositionNotFound,
		},
		{
			name: "Non-existing employee",
			request: domain.UpdateEmployee{
				EmployeeID: uuid.New(),
				PositionID: positionID,
				FirstName:  "NewName",
				LastName:   "NewLastName",
			},
			wantErr: customerrors.ErrEmployeeNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := repo.UpdateEmployee(ctx, tt.request)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestEmployeeRepository_DeleteEmployee(t *testing.T) {
	ctx := context.Background()

	repo, container := setupEmployeeRepo(ctx, t)
	defer func(container *postgres.PostgresContainer, ctx context.Context) {
		err := container.Terminate(ctx)
		if err != nil {
			t.Fatalf("could not terminate postgres container: %v", err.Error())
		}
	}(container, ctx)

	positionID := uuid.New()
	employeeID := uuid.New()

	_, err := repo.CreateEmployee(ctx, domain.CreateEmployee{
		EmployeeID:   employeeID,
		PositionID:   positionID,
		FirstName:    "John",
		LastName:     "Doe",
		PositionName: "Go developer",
		Salary:       30999,
	})

	require.NoError(t, err)

	tests := []struct {
		name       string
		employeeID uuid.UUID
		wantErr    error
	}{
		{
			name:       "Valid employee id",
			employeeID: employeeID,
		},
		{
			name:       "Non-existing employee",
			employeeID: uuid.New(),
			wantErr:    customerrors.ErrEmployeeNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.DeleteEmployee(ctx, tt.employeeID)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}