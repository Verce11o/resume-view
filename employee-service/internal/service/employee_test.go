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
	t.Parallel()

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
		{
			name: "Invalid employee ID",
			input: domain.CreateEmployee{
				EmployeeID:   uuid.Nil,
				PositionID:   positionID,
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
					}, nil)

				f.employeeRepo.On("CreateEmployee", mock.Anything, mock.AnythingOfType("domain.CreateEmployee")).
					Return(models.Employee{
						ID:         employeeID,
						FirstName:  "John",
						LastName:   "Doe",
						PositionID: positionID,
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

			assert.Equal(t, tt.wantErr, err != nil)

			assert.EqualValues(t, tt.response, employee)
		})
	}
}

func TestEmployeeService_GetEmployee(t *testing.T) {
	t.Parallel()

	type fields struct {
		employeeRepo *mocks.EmployeeRepository
		positionRepo *mocks.PositionRepository
		cache        *mocks.EmployeeCacheRepository
	}

	employeeID := uuid.New()

	tests := []struct {
		name     string
		id       uuid.UUID
		response models.Employee
		mockFunc func(f *fields)
		wantErr  bool
	}{
		{
			name: "Valid return from cache",
			id:   employeeID,
			response: models.Employee{
				ID:        employeeID,
				FirstName: "John",
				LastName:  "Doe",
			},
			mockFunc: func(f *fields) {
				f.cache.On("GetEmployee", mock.Anything, mock.AnythingOfType("string")).
					Return(&models.Employee{
						ID:        employeeID,
						FirstName: "John",
						LastName:  "Doe",
					}, nil)
			},
		},
		{
			name: "Valid return from repo",
			id:   employeeID,
			response: models.Employee{
				ID:        employeeID,
				FirstName: "John",
				LastName:  "Doe",
			},
			mockFunc: func(f *fields) {
				f.cache.On("GetEmployee", mock.Anything, mock.AnythingOfType("string")).
					Return(nil, assert.AnError)

				f.employeeRepo.On("GetEmployee", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(models.Employee{
						ID:        employeeID,
						FirstName: "John",
						LastName:  "Doe",
					}, nil)

				f.cache.On("SetEmployee", mock.Anything, mock.AnythingOfType("string"),
					mock.AnythingOfType("*models.Employee")).
					Return(nil)
			},
		},
		{
			name:     "Invalid employee ID",
			id:       uuid.Nil,
			response: models.Employee{},
			mockFunc: func(f *fields) {
				f.cache.On("GetEmployee", mock.Anything, mock.AnythingOfType("string")).
					Return(nil, assert.AnError)

				f.employeeRepo.On("GetEmployee", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(models.Employee{}, assert.AnError)
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
				cache:        cache,
			})

			srv := &EmployeeService{
				log:          zap.NewNop().Sugar(),
				employeeRepo: employeeRepo,
				positionRepo: positionRepo,
				cache:        cache,
				transactor:   transactor,
			}

			employee, err := srv.GetEmployee(context.TODO(), tt.id)

			assert.Equal(t, tt.wantErr, err != nil)

			assert.EqualValues(t, tt.response, employee)
		})
	}
}

func TestEmployeeService_GetEmployeeList(t *testing.T) {
	t.Parallel()

	type fields struct {
		employeeRepo *mocks.EmployeeRepository
		positionRepo *mocks.PositionRepository
		cache        *mocks.EmployeeCacheRepository
	}

	employeeID := uuid.New()

	tests := []struct {
		name     string
		cursor   string
		response models.EmployeeList
		mockFunc func(f *fields)
		wantErr  bool
	}{
		{
			name:   "Valid empty cursor",
			cursor: "",
			response: models.EmployeeList{
				Cursor: "cursorExample",
				Employees: []models.Employee{
					{
						ID:        employeeID,
						FirstName: "John",
						LastName:  "Doe",
					},
				},
			},
			mockFunc: func(f *fields) {
				f.employeeRepo.On("GetEmployeeList", mock.Anything, mock.AnythingOfType("string")).
					Return(models.EmployeeList{
						Cursor: "cursorExample",
						Employees: []models.Employee{
							{
								ID:        employeeID,
								FirstName: "John",
								LastName:  "Doe",
							},
						},
					}, nil)
			},
		},
		{
			name:     "Invalid cursor",
			cursor:   "invalid",
			response: models.EmployeeList{},
			mockFunc: func(f *fields) {
				f.employeeRepo.On("GetEmployeeList", mock.Anything, mock.AnythingOfType("string")).
					Return(models.EmployeeList{}, assert.AnError)
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
				cache:        cache,
			})

			srv := &EmployeeService{
				log:          zap.NewNop().Sugar(),
				employeeRepo: employeeRepo,
				positionRepo: positionRepo,
				cache:        cache,
				transactor:   transactor,
			}

			employee, err := srv.GetEmployeeList(context.TODO(), tt.cursor)

			assert.Equal(t, tt.wantErr, err != nil)

			assert.EqualValues(t, tt.response, employee)
		})
	}
}

func TestEmployeeService_UpdateEmployee(t *testing.T) {
	t.Parallel()

	type fields struct {
		employeeRepo *mocks.EmployeeRepository
		positionRepo *mocks.PositionRepository
		cache        *mocks.EmployeeCacheRepository
	}

	employeeID := uuid.New()
	positionID := uuid.New()

	tests := []struct {
		name     string
		input    domain.UpdateEmployee
		response models.Employee
		mockFunc func(f *fields)
		wantErr  bool
	}{
		{
			name: "Valid input",
			input: domain.UpdateEmployee{
				EmployeeID: employeeID,
				PositionID: positionID,
				FirstName:  "John",
				LastName:   "Doe",
				Salary:     30999,
			},
			response: models.Employee{
				ID:         employeeID,
				FirstName:  "John",
				LastName:   "Doe",
				PositionID: positionID,
			},
			mockFunc: func(f *fields) {
				f.employeeRepo.On("UpdateEmployee", mock.Anything, mock.AnythingOfType("domain.UpdateEmployee")).
					Return(models.Employee{
						ID:         employeeID,
						FirstName:  "John",
						LastName:   "Doe",
						PositionID: positionID,
					}, nil)

				f.cache.On("DeleteEmployee", mock.Anything, mock.AnythingOfType("string")).
					Return(nil)
			},
		},
		{
			name: "Invalid input",
			input: domain.UpdateEmployee{
				EmployeeID: employeeID,
				PositionID: uuid.Nil,
				FirstName:  "John",
				LastName:   "Doe",
				Salary:     30999,
			},
			response: models.Employee{},
			mockFunc: func(f *fields) {
				f.employeeRepo.On("UpdateEmployee", mock.Anything, mock.AnythingOfType("domain.UpdateEmployee")).
					Return(models.Employee{}, assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "Cache error",
			input: domain.UpdateEmployee{
				EmployeeID: employeeID,
				PositionID: positionID,
				FirstName:  "John",
				LastName:   "Doe",
				Salary:     30999,
			},
			response: models.Employee{
				ID:         employeeID,
				FirstName:  "John",
				LastName:   "Doe",
				PositionID: positionID,
			},
			mockFunc: func(f *fields) {
				f.employeeRepo.On("UpdateEmployee", mock.Anything, mock.AnythingOfType("domain.UpdateEmployee")).
					Return(models.Employee{
						ID:         employeeID,
						FirstName:  "John",
						LastName:   "Doe",
						PositionID: positionID,
					}, nil)

				f.cache.On("DeleteEmployee", mock.Anything, mock.AnythingOfType("string")).
					Return(assert.AnError)
			},
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
				cache:        cache,
			})

			srv := &EmployeeService{
				log:          zap.NewNop().Sugar(),
				employeeRepo: employeeRepo,
				positionRepo: positionRepo,
				cache:        cache,
				transactor:   transactor,
			}
			employee, err := srv.UpdateEmployee(context.TODO(), tt.input)

			assert.Equal(t, tt.wantErr, err != nil)

			assert.EqualValues(t, tt.response, employee)
		})
	}
}

func TestEmployeeService_DeleteEmployee(t *testing.T) {
	t.Parallel()

	type fields struct {
		employeeRepo *mocks.EmployeeRepository
		positionRepo *mocks.PositionRepository
		cache        *mocks.EmployeeCacheRepository
	}

	employeeID := uuid.New()
	positionID := uuid.New()

	tests := []struct {
		name     string
		id       uuid.UUID
		mockFunc func(f *fields)
		wantErr  bool
	}{
		{
			name: "Valid input",
			id:   employeeID,
			mockFunc: func(f *fields) {
				f.employeeRepo.On("GetEmployee", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(models.Employee{
						ID:         employeeID,
						FirstName:  "John",
						LastName:   "Doe",
						PositionID: positionID,
					}, nil)

				f.employeeRepo.On("DeleteEmployee", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(nil)

				f.cache.On("DeleteEmployee", mock.Anything, mock.AnythingOfType("string")).
					Return(nil)
			},
		},
		{
			name: "Non-existing employee",
			id:   uuid.New(),
			mockFunc: func(f *fields) {
				f.employeeRepo.On("GetEmployee", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(models.Employee{}, assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "Repository error",
			id:   employeeID,
			mockFunc: func(f *fields) {
				f.employeeRepo.On("GetEmployee", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(models.Employee{
						ID:        employeeID,
						FirstName: "John",
						LastName:  "Doe",
					}, nil)

				f.employeeRepo.On("DeleteEmployee", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "Cache error",
			id:   employeeID,
			mockFunc: func(f *fields) {
				f.employeeRepo.On("GetEmployee", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(models.Employee{
						ID:        employeeID,
						FirstName: "John",
						LastName:  "Doe",
					}, nil)

				f.employeeRepo.On("DeleteEmployee", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(nil)

				f.cache.On("DeleteEmployee", mock.Anything, mock.AnythingOfType("string")).
					Return(assert.AnError)
			},
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
				cache:        cache,
			})

			srv := &EmployeeService{
				log:          zap.NewNop().Sugar(),
				employeeRepo: employeeRepo,
				positionRepo: positionRepo,
				cache:        cache,
				transactor:   transactor,
			}

			err := srv.DeleteEmployee(context.TODO(), tt.id)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
