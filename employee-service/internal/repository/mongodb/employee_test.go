//go:build integration

package mongodb

import (
	"context"
	"testing"
	"time"

	"github.com/Verce11o/resume-view/employee-service/internal/domain"
	"github.com/Verce11o/resume-view/employee-service/internal/lib/customerrors"
	"github.com/Verce11o/resume-view/employee-service/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type EmployeeRepositorySuite struct {
	suite.Suite
	positionID uuid.UUID
	ctx        context.Context
	client     *mongo.Client
	container  testcontainers.Container
	repo       *EmployeeRepository
}

func (s *EmployeeRepositorySuite) SetupSuite() {
	s.ctx = context.Background()
	container, connURI := setupMongoDBContainer(s.ctx, s.T())

	client, err := mongo.Connect(s.ctx,
		options.Client().ApplyURI(connURI),
		options.Client().SetMaxConnIdleTime(3*time.Second))
	require.NoError(s.T(), err)

	s.repo = NewEmployeeRepository(client.Database("employees"))
	positionRepo := NewPositionRepository(client.Database("employees"))

	s.positionID = uuid.New()
	_, err = positionRepo.CreatePosition(s.ctx, domain.CreatePosition{
		ID:     s.positionID,
		Name:   "Go Developer",
		Salary: 10999,
	})

	s.client = client
	s.container = container
	require.NoError(s.T(), err)
}

func (s *EmployeeRepositorySuite) TearDownSuite() {
	err := s.container.Terminate(s.ctx)
	require.NoError(s.T(), err)
}

func (s *EmployeeRepositorySuite) TestCreateEmployee() {

	employeeID := uuid.New()

	tests := []struct {
		name     string
		request  domain.CreateEmployee
		response models.Employee
		wantErr  error
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
			response: models.Employee{
				ID:         employeeID,
				PositionID: s.positionID,
				FirstName:  "John",
				LastName:   "Doe",
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
			response: models.Employee{},
			wantErr:  customerrors.ErrDuplicateID,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			resp, err := s.repo.CreateEmployee(s.ctx, tt.request)
			assert.EqualExportedValues(s.T(), tt.response, resp)
			assert.ErrorIs(s.T(), err, tt.wantErr)
		})
	}
}

func (s *EmployeeRepositorySuite) TestGetEmployee() {

	employeeID := uuid.New()

	employee, err := s.repo.CreateEmployee(s.ctx, domain.CreateEmployee{
		EmployeeID:   employeeID,
		PositionID:   s.positionID,
		FirstName:    "John",
		LastName:     "Doe",
		PositionName: "C++ Developer",
		Salary:       30999,
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
			name:       "Non-existent employee id",
			employeeID: uuid.Nil,
			wantErr:    customerrors.ErrEmployeeNotFound,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {

			resp, err := s.repo.GetEmployee(s.ctx, tt.employeeID)
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
		_, err := s.repo.CreateEmployee(s.ctx, domain.CreateEmployee{
			EmployeeID:   uuid.New(),
			PositionID:   s.positionID,
			FirstName:    "Sample",
			LastName:     "Sample",
			PositionName: "Python Developer",
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

			resp, err := s.repo.GetEmployeeList(s.ctx, tt.cursor)
			assert.ErrorIs(s.T(), err, tt.wantErr)
			assert.Equal(s.T(), len(resp.Employees), tt.length)
			nextCursor = resp.Cursor
		})
	}
}

func (s *EmployeeRepositorySuite) TestUpdateEmployee() {

	employeeID := uuid.New()

	_, err := s.repo.CreateEmployee(s.ctx, domain.CreateEmployee{
		EmployeeID:   employeeID,
		PositionID:   s.positionID,
		FirstName:    "Sample",
		LastName:     "Sample",
		PositionName: "Python Developer",
		Salary:       30999,
	})
	require.NoError(s.T(), err)

	tests := []struct {
		name     string
		request  domain.UpdateEmployee
		response models.Employee
		wantErr  error
	}{
		{
			name: "Valid input",
			request: domain.UpdateEmployee{
				EmployeeID: employeeID,
				PositionID: s.positionID,
				FirstName:  "New Name",
				LastName:   "New Last Name",
			},
			response: models.Employee{
				ID:         employeeID,
				PositionID: s.positionID,
				FirstName:  "New Name",
				LastName:   "New Last Name",
			},
		},
		{
			name: "Non-existent employee id",
			request: domain.UpdateEmployee{
				EmployeeID: uuid.Nil,
				PositionID: s.positionID,
				FirstName:  "New Name",
				LastName:   "New Last Name",
			},
			response: models.Employee{},
			wantErr:  customerrors.ErrEmployeeNotFound,
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {

			resp, err := s.repo.UpdateEmployee(s.ctx, tt.request)
			assert.ErrorIs(s.T(), err, tt.wantErr)
			assert.EqualExportedValues(s.T(), tt.response, resp)
		})
	}
}

func (s *EmployeeRepositorySuite) TestDeleteEmployee() {

	employeeID := uuid.New()

	_, err := s.repo.CreateEmployee(s.ctx, domain.CreateEmployee{
		EmployeeID:   employeeID,
		PositionID:   s.positionID,
		FirstName:    "John",
		LastName:     "Doe",
		PositionName: "C++ Developer",
		Salary:       30999,
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
			name:       "Non-existent employee id",
			employeeID: uuid.Nil,
			wantErr:    customerrors.ErrEmployeeNotFound,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			err = s.repo.DeleteEmployee(s.ctx, tt.employeeID)
			assert.ErrorIs(s.T(), err, tt.wantErr)

			if tt.wantErr != nil {
				_, err = s.repo.GetEmployee(s.ctx, employeeID)
				assert.ErrorIs(s.T(), err, customerrors.ErrEmployeeNotFound)
			}
		})
	}
}

func TestEmployeeRepositorySuite(t *testing.T) {
	suite.Run(t, new(EmployeeRepositorySuite))
}
