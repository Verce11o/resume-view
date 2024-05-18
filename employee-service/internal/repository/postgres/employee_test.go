//go:build integration

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
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

type EmployeeRepositorySuite struct {
	suite.Suite
	ctx          context.Context
	positionID   uuid.UUID
	employeeRepo *EmployeeRepository
	container    *postgres.PostgresContainer
}

func (s *EmployeeRepositorySuite) SetupSuite() {
	s.ctx = context.Background()

	container, connURI := setupPostgresContainer(s.ctx, s.T())
	dbPool, err := pgxpool.New(s.ctx, connURI)
	require.NoError(s.T(), err)

	employeeRepo := NewEmployeeRepository(dbPool)
	positionRepo := NewPositionRepository(dbPool)

	positionID := uuid.New()
	s.positionID = positionID

	_, err = positionRepo.CreatePosition(s.ctx, domain.CreatePosition{
		ID:     positionID,
		Name:   "Go Developer",
		Salary: 10999,
	})

	require.NoError(s.T(), err)

	s.employeeRepo = employeeRepo
	s.container = container
}

func (s *EmployeeRepositorySuite) TearDownSuite() {
	err := s.container.Terminate(s.ctx)
	if err != nil {
		s.T().Fatalf("could not terminate postgres container: %v", err.Error())
	}
}

func (s *EmployeeRepositorySuite) TestCreateEmployee() {
	employeeID := uuid.New()

	tests := []struct {
		name    string
		request domain.CreateEmployee
		wantErr error
	}{
		{
			name: "Valid input",
			request: domain.CreateEmployee{
				EmployeeID:   employeeID,
				PositionID:   s.positionID,
				FirstName:    "John",
				LastName:     "Doe",
				PositionName: "Go Developer",
				Salary:       12345,
			},
		},
		{
			name: "Duplicate employee id",
			request: domain.CreateEmployee{
				EmployeeID:   employeeID,
				PositionID:   s.positionID,
				FirstName:    "John",
				LastName:     "Doe",
				PositionName: "Python Developer",
				Salary:       12345,
			},
			wantErr: customerrors.ErrDuplicateID,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			_, err := s.employeeRepo.CreateEmployee(s.ctx, tt.request)
			assert.ErrorIs(s.T(), err, tt.wantErr)
		})
	}
}

func (s *EmployeeRepositorySuite) TestGetEmployee() {
	employeeID := uuid.New()

	employee, err := s.employeeRepo.CreateEmployee(s.ctx, domain.CreateEmployee{
		EmployeeID:   employeeID,
		PositionID:   s.positionID,
		FirstName:    "John",
		LastName:     "Doe",
		PositionName: "Go Developer",
		Salary:       0,
	})

	require.NoError(s.T(), err)

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
		s.Run(tt.name, func() {
			resp, err := s.employeeRepo.GetEmployee(s.ctx, tt.employeeID)
			if tt.wantErr != nil {
				assert.ErrorIs(s.T(), err, tt.wantErr)
				assert.NotEqual(s.T(), resp, employee)

				return
			}

			assert.Equal(s.T(), resp, employee)
		})
	}
}

func (s *EmployeeRepositorySuite) TestGetEmployeeList() {
	for i := 0; i < 10; i++ {
		_, err := s.employeeRepo.CreateEmployee(s.ctx, domain.CreateEmployee{
			EmployeeID:   uuid.New(),
			PositionID:   s.positionID,
			FirstName:    "John",
			LastName:     "Doe",
			PositionName: "Go Developer",
			Salary:       30999,
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
			resp, err := s.employeeRepo.GetEmployeeList(s.ctx, tt.cursor)
			assert.ErrorIs(s.T(), err, tt.wantErr)
			assert.Equal(s.T(), len(resp.Employees), tt.length)
			nextCursor = resp.Cursor
		})
	}
}

func (s *EmployeeRepositorySuite) TestUpdateEmployee() {
	employeeID := uuid.New()

	_, err := s.employeeRepo.CreateEmployee(s.ctx, domain.CreateEmployee{
		EmployeeID:   employeeID,
		PositionID:   s.positionID,
		FirstName:    "John",
		LastName:     "Doe",
		PositionName: "Go developer",
		Salary:       30999,
	})
	require.NoError(s.T(), err)

	tests := []struct {
		name    string
		request domain.UpdateEmployee
		wantErr error
	}{
		{
			name: "Valid input",
			request: domain.UpdateEmployee{
				EmployeeID: employeeID,
				PositionID: s.positionID,
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
				PositionID: s.positionID,
				FirstName:  "NewName",
				LastName:   "NewLastName",
			},
			wantErr: customerrors.ErrEmployeeNotFound,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			_, err := s.employeeRepo.UpdateEmployee(s.ctx, tt.request)
			assert.ErrorIs(s.T(), err, tt.wantErr)
		})
	}
}

func (s *EmployeeRepositorySuite) TestDeleteEmployee() {
	employeeID := uuid.New()

	_, err := s.employeeRepo.CreateEmployee(s.ctx, domain.CreateEmployee{
		EmployeeID:   employeeID,
		PositionID:   s.positionID,
		FirstName:    "John",
		LastName:     "Doe",
		PositionName: "Go developer",
		Salary:       30999,
	})
	require.NoError(s.T(), err)

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
		s.Run(tt.name, func() {
			err := s.employeeRepo.DeleteEmployee(s.ctx, tt.employeeID)
			assert.ErrorIs(s.T(), err, tt.wantErr)
		})
	}
}

func TestEmployeeRepositorySuite(t *testing.T) {
	suite.Run(t, new(EmployeeRepositorySuite))
}
