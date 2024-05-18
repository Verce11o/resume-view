//go:build integration

package redis

import (
	"context"
	"testing"
	"time"

	"github.com/Verce11o/resume-view/employee-service/internal/lib/customerrors"
	"github.com/Verce11o/resume-view/employee-service/internal/models"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	redisContainer "github.com/testcontainers/testcontainers-go/modules/redis"
)

type EmployeeCacheSuite struct {
	suite.Suite
	ctx       context.Context
	client    *redis.Client
	container *redisContainer.RedisContainer
	repo      *EmployeeCache
}

func (s *EmployeeCacheSuite) SetupSuite() {
	s.ctx = context.Background()
	container, connURI := setupRedisContainer(s.ctx, s.T())

	client := redis.NewClient(&redis.Options{
		Addr: connURI,
	})

	s.client = client
	s.repo = NewEmployeeCache(client)
	s.container = container
}

func (s *EmployeeCacheSuite) TearDownSuite() {
	err := s.container.Terminate(s.ctx)
	require.NoError(s.T(), err)

	err = s.client.Close()
	require.NoError(s.T(), err)
}

func (s *EmployeeCacheSuite) TestSetEmployee() {

	employeeID := uuid.New().String()

	tests := []struct {
		name       string
		employeeID string
		employee   *models.Employee
		wantErr    bool
	}{
		{
			name:       "Valid input",
			employeeID: employeeID,
			employee: &models.Employee{
				ID:         uuid.New(),
				FirstName:  "John",
				LastName:   "Doe",
				PositionID: uuid.New(),
				CreatedAt:  time.Now().UTC(),
				UpdatedAt:  time.Now().UTC(),
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {

			err := s.repo.SetEmployee(s.ctx, tt.employeeID, tt.employee)
			if tt.wantErr {
				assert.Error(s.T(), err)
			} else {
				assert.NoError(s.T(), err)
			}
		})
	}
}

func (s *EmployeeCacheSuite) TestGetEmployee() {

	employeeID := uuid.New()

	err := s.repo.SetEmployee(s.ctx, employeeID.String(), &models.Employee{
		ID:        employeeID,
		FirstName: "John",
		LastName:  "Doe",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})

	require.NoError(s.T(), err)

	tests := []struct {
		name       string
		employeeID string
		wantErr    error
	}{
		{
			name:       "Valid input",
			employeeID: employeeID.String(),
		},
		{
			name:       "Non-existent employee",
			employeeID: uuid.New().String(),
			wantErr:    customerrors.ErrEmployeeNotCached,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {

			_, err := s.repo.GetEmployee(s.ctx, tt.employeeID)
			assert.ErrorIs(s.T(), err, tt.wantErr)
		})
	}
}

func (s *EmployeeCacheSuite) TestDeleteEmployee() {

	employeeID := uuid.New()

	err := s.repo.SetEmployee(s.ctx, employeeID.String(), &models.Employee{
		ID:         employeeID,
		FirstName:  "John",
		LastName:   "Doe",
		PositionID: uuid.New(),
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	})

	require.NoError(s.T(), err)

	tests := []struct {
		name       string
		employeeID string
		wantErr    error
	}{
		{
			name:       "Valid input",
			employeeID: employeeID.String(),
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {

			err := s.repo.DeleteEmployee(s.ctx, tt.employeeID)
			assert.ErrorIs(s.T(), err, tt.wantErr)
		})
	}
}

func TestEmployeeCacheSuite(t *testing.T) {
	suite.Run(t, new(EmployeeCacheSuite))
}
