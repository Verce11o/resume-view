package mongodb

import (
	"context"
	"github.com/Verce11o/resume-view/employee-service/api"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
	"time"
)

func TestEmployeeRepository_CreateEmployee(t *testing.T) {
	ctx := context.Background()

	container, connURI := setupMongoDBContainer(t)
	defer container.Terminate(ctx)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connURI), options.Client().SetMaxConnIdleTime(3*time.Second))

	require.NoError(t, err)

	repo := NewEmployeeRepository(client.Database("employees"))

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
			}
		})
	}

}

func TestEmployeeRepository_GetEmployee(t *testing.T) {
	ctx := context.Background()

	container, connURI := setupMongoDBContainer(t)
	defer container.Terminate(ctx)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connURI), options.Client().SetMaxConnIdleTime(3*time.Second))

	require.NoError(t, err)

	err = client.Ping(ctx, nil)

	require.NoError(t, err)

	repo := NewEmployeeRepository(client.Database("employee"))

	employeeID := uuid.New()
	positionID := uuid.New()

	_, err = repo.CreateEmployee(ctx, employeeID, positionID, api.CreateEmployee{
		FirstName:    "John",
		LastName:     "Doe",
		PositionName: "C++ Developer",
		Salary:       30999,
	})

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
			name:       "Non-existent employee id",
			employeeID: uuid.Nil,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := repo.GetEmployee(ctx, tt.employeeID)
			if tt.wantErr {
				assert.Error(t, err)
			}
		})
	}
}

func TestEmployeeRepository_GetEmployeeList(t *testing.T) {
	ctx := context.Background()

	container, connURI := setupMongoDBContainer(t)
	defer container.Terminate(ctx)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connURI), options.Client().SetMaxConnIdleTime(3*time.Second))

	require.NoError(t, err)

	err = client.Ping(ctx, nil)

	require.NoError(t, err)

	repo := NewEmployeeRepository(client.Database("employee"))

	for i := 0; i < 10; i++ {
		_, err = repo.CreateEmployee(ctx, uuid.New(), uuid.New(), api.CreateEmployee{
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
			}
			assert.Equal(t, len(resp.Employees), tt.length)
			nextCursor = resp.Cursor
		})
	}
}

func TestEmployeeRepository_UpdateEmployee(t *testing.T) {
	ctx := context.Background()

	container, connURI := setupMongoDBContainer(t)
	defer container.Terminate(ctx)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connURI), options.Client().SetMaxConnIdleTime(3*time.Second))

	require.NoError(t, err)

	err = client.Ping(ctx, nil)

	require.NoError(t, err)

	repo := NewEmployeeRepository(client.Database("employee"))

	employeeID := uuid.New()
	positionID := uuid.New()

	_, err = repo.CreateEmployee(ctx, employeeID, positionID, api.CreateEmployee{
		FirstName:    "Sample",
		LastName:     "Sample",
		PositionName: "Python Developer",
		Salary:       30999,
	})

	require.NoError(t, err)

	tests := []struct {
		name       string
		employeeID uuid.UUID
		request    api.UpdateEmployee
		wantErr    bool
	}{
		{
			name:       "Valid input",
			employeeID: employeeID,
			request: api.UpdateEmployee{
				FirstName:  "New Name",
				LastName:   "New Last Name",
				PositionId: positionID.String(),
			},
		},
		{
			name:       "Non-existent employee id",
			employeeID: uuid.Nil,
			request: api.UpdateEmployee{
				FirstName:  "New Name",
				LastName:   "New Last  Name",
				PositionId: positionID.String(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err = repo.UpdateEmployee(ctx, tt.employeeID, tt.request)
			if tt.wantErr {
				assert.Error(t, err)
			}

		})
	}

}

func TestEmployeeRepository_DeleteEmployee(t *testing.T) {
	ctx := context.Background()

	container, connURI := setupMongoDBContainer(t)
	defer container.Terminate(ctx)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connURI), options.Client().SetMaxConnIdleTime(3*time.Second))

	require.NoError(t, err)

	err = client.Ping(ctx, nil)

	require.NoError(t, err)

	repo := NewEmployeeRepository(client.Database("employee"))

	employeeID := uuid.New()
	positionID := uuid.New()

	_, err = repo.CreateEmployee(ctx, employeeID, positionID, api.CreateEmployee{
		FirstName:    "John",
		LastName:     "Doe",
		PositionName: "C++ Developer",
		Salary:       30999,
	})
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
			name:       "Non-existent employee id",
			employeeID: uuid.Nil,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err = repo.DeleteEmployee(ctx, tt.employeeID)
			if tt.wantErr {
				assert.Error(t, err)
			}
		})
	}
}
