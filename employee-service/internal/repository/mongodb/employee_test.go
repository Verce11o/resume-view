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
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setupEmployeeRepo(ctx context.Context, t *testing.T) (*EmployeeRepository, testcontainers.Container) {
	container, connURI := setupMongoDBContainer(ctx, t)

	client, err := mongo.Connect(ctx,
		options.Client().ApplyURI(connURI),
		options.Client().SetMaxConnIdleTime(3*time.Second))

	require.NoError(t, err)

	repo := NewEmployeeRepository(client.Database("employees"))

	return repo, container
}

func TestEmployeeRepository_CreateEmployee(t *testing.T) {
	ctx := context.Background()

	repo, container := setupEmployeeRepo(ctx, t)
	defer func(container testcontainers.Container, ctx context.Context) {
		err := container.Terminate(ctx)
		if err != nil {
			t.Fatalf("could not terminate mongo container: %v", err.Error())
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
	defer func(container testcontainers.Container, ctx context.Context) {
		err := container.Terminate(ctx)
		if err != nil {
			t.Fatalf("could not terminate mongo container: %v", err.Error())
		}
	}(container, ctx)

	employeeID := uuid.New()
	positionID := uuid.New()

	employee, err := repo.CreateEmployee(ctx, domain.CreateEmployee{
		EmployeeID:   employeeID,
		PositionID:   positionID,
		FirstName:    "John",
		LastName:     "Doe",
		PositionName: "C++ Developer",
		Salary:       30999,
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
			name:       "Non-existent employee id",
			employeeID: uuid.Nil,
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
	defer func(container testcontainers.Container, ctx context.Context) {
		err := container.Terminate(ctx)
		if err != nil {
			t.Fatalf("could not terminate mongo container: %v", err.Error())
		}
	}(container, ctx)

	for i := 0; i < 10; i++ {
		_, err := repo.CreateEmployee(ctx, domain.CreateEmployee{
			EmployeeID:   uuid.New(),
			PositionID:   uuid.New(),
			FirstName:    "Sample",
			LastName:     "Sample",
			PositionName: "Python Developer",
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
	defer func(container testcontainers.Container, ctx context.Context) {
		err := container.Terminate(ctx)
		if err != nil {
			t.Fatalf("could not terminate mongo container: %v", err.Error())
		}
	}(container, ctx)

	employeeID := uuid.New()
	positionID := uuid.New()

	_, err := repo.CreateEmployee(ctx, domain.CreateEmployee{
		EmployeeID:   employeeID,
		PositionID:   positionID,
		FirstName:    "Sample",
		LastName:     "Sample",
		PositionName: "Python Developer",
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
				FirstName:  "New Name",
				LastName:   "New Last Name",
			},
		},
		{
			name: "Non-existent employee id",
			request: domain.UpdateEmployee{
				EmployeeID: uuid.Nil,
				PositionID: positionID,
				FirstName:  "New Name",
				LastName:   "New Last  Name",
			},
			wantErr: customerrors.ErrEmployeeNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err = repo.UpdateEmployee(ctx, tt.request)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestEmployeeRepository_DeleteEmployee(t *testing.T) {
	ctx := context.Background()

	repo, container := setupEmployeeRepo(ctx, t)
	defer func(container testcontainers.Container, ctx context.Context) {
		err := container.Terminate(ctx)
		if err != nil {
			t.Fatalf("could not terminate mongo container: %v", err.Error())
		}
	}(container, ctx)

	employeeID := uuid.New()
	positionID := uuid.New()

	_, err := repo.CreateEmployee(ctx, domain.CreateEmployee{
		EmployeeID:   employeeID,
		PositionID:   positionID,
		FirstName:    "John",
		LastName:     "Doe",
		PositionName: "C++ Developer",
		Salary:       30999,
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
			name:       "Non-existent employee id",
			employeeID: uuid.Nil,
			wantErr:    customerrors.ErrEmployeeNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err = repo.DeleteEmployee(ctx, tt.employeeID)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
