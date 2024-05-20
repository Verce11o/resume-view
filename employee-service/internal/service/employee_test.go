package service

import (
	"context"
	"testing"

	"github.com/Verce11o/resume-view/employee-service/internal/domain"
	"github.com/Verce11o/resume-view/employee-service/internal/models"
	"github.com/Verce11o/resume-view/employee-service/internal/service/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestEmployeeService_CreateEmployee(t *testing.T) {
	type fields struct {
		employeeRepo *mocks.EmployeeRepository
		positionRepo *mocks.PositionRepository
		transactor   *mocks.Transactor
	}

	employeeID := uuid.New()
	positionID := uuid.New()

	tests := []struct {
		name     string
		input    domain.CreateEmployee
		response models.Employee
		mockFunc func(f *fields)
		wantErr  bool
	}{
		{
			name: "Valid",
			input: domain.CreateEmployee{
				EmployeeID:   employeeID,
				PositionID:   positionID,
				FirstName:    "John",
				LastName:     "Doe",
				PositionName: "Go Developer",
				Salary:       30999,
			},
			response: models.Employee{
				ID:         employeeID,
				FirstName:  "John",
				LastName:   "Doe",
				PositionID: positionID,
			},
			mockFunc: func(f *fields) {
				f.transactor.On("WithTransaction", mock.Anything, mock.Anything).
					Return(nil).Run(func(args mock.Arguments) {
					fn := args.Get(1).(func(ctx context.Context) error)
					assert.NoError(t, fn(context.TODO()))
				})

				f.positionRepo.On("CreatePosition", mock.Anything, mock.AnythingOfType("domain.CreatePosition")).
					Return(models.Position{
						ID:     positionID,
						Name:   "Go Developer",
						Salary: 30999,
					}, nil)

				f.employeeRepo.On("CreateEmployee", mock.Anything, mock.AnythingOfType("domain.CreateEmployee")).
					Return(models.Employee{
						ID:         employeeID,
						FirstName:  "John",
						LastName:   "Doe",
						PositionID: positionID,
					}, nil)
			},
		},
		{
			name: "Invalid position ID",
			input: domain.CreateEmployee{
				EmployeeID:   employeeID,
				PositionID:   uuid.Nil,
				FirstName:    "John",
				LastName:     "Doe",
				PositionName: "Go Developer",
				Salary:       30999,
			},
			response: models.Employee{},
			mockFunc: func(f *fields) {
				f.transactor.On("WithTransaction", mock.Anything, mock.Anything).
					Return(assert.AnError).Run(func(args mock.Arguments) {
					fn := args.Get(1).(func(ctx context.Context) error)
					assert.Error(t, fn(context.TODO()))
				})

				f.positionRepo.On("CreatePosition", mock.Anything, mock.AnythingOfType("domain.CreatePosition")).
					Return(models.Position{
						ID:     positionID,
						Name:   "Go Developer",
						Salary: 30999,
					}, assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			employeeRepo := mocks.NewEmployeeRepository(t)
			positionRepo := mocks.NewPositionRepository(t)
			transactor := mocks.NewTransactor(t)
			cache := mocks.NewEmployeeCacheRepository(t)

			tt.mockFunc(&fields{
				employeeRepo: employeeRepo,
				positionRepo: positionRepo,
				transactor:   transactor,
			})

			srv := &EmployeeService{
				log:          zap.NewNop().Sugar(),
				employeeRepo: employeeRepo,
				positionRepo: positionRepo,
				cache:        cache,
				transactor:   transactor,
			}

			employee, err := srv.CreateEmployee(context.TODO(), tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			}

			assert.EqualValues(t, tt.response, employee)
		})
	}
}
