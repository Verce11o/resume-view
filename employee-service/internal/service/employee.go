package service

import (
	"context"
	"fmt"

	"github.com/Verce11o/resume-view/employee-service/internal/domain"
	"github.com/Verce11o/resume-view/employee-service/internal/models"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

//go:generate go run github.com/vektra/mockery/v2@v2.42.2 --name=EmployeeRepository
type EmployeeRepository interface {
	CreateEmployee(ctx context.Context, req domain.CreateEmployee) (models.Employee, error)
	GetEmployee(ctx context.Context, id uuid.UUID) (models.Employee, error)
	GetEmployeeList(ctx context.Context, cursor string) (models.EmployeeList, error)
	UpdateEmployee(ctx context.Context, req domain.UpdateEmployee) (models.Employee, error)
	DeleteEmployee(ctx context.Context, id uuid.UUID) error
}

//go:generate go run github.com/vektra/mockery/v2@v2.42.2 --name=EmployeeCacheRepository
type EmployeeCacheRepository interface {
	GetEmployee(ctx context.Context, key string) (*models.Employee, error)
	SetEmployee(ctx context.Context, employeeID string, employee *models.Employee) error
	DeleteEmployee(ctx context.Context, employeeID string) error
}

//go:generate go run github.com/vektra/mockery/v2@v2.42.2 --name=Transactor
type Transactor interface {
	WithTransaction(context.Context, func(ctx context.Context) error) error
}

type EventNotifier interface {
	SendMessage(ctx context.Context, key, value []byte) error
}

type EmployeeService struct {
	log           *zap.SugaredLogger
	employeeRepo  EmployeeRepository
	positionRepo  PositionRepository
	cache         EmployeeCacheRepository
	transactor    Transactor
	eventNotifier EventNotifier
}

func NewEmployeeService(log *zap.SugaredLogger, employeeRepo EmployeeRepository, positionRepo PositionRepository,
	cache EmployeeCacheRepository, transactor Transactor, notifier EventNotifier) *EmployeeService {
	return &EmployeeService{log: log, employeeRepo: employeeRepo, positionRepo: positionRepo, cache: cache,
		transactor: transactor, eventNotifier: notifier}
}

func (s *EmployeeService) CreateEmployee(ctx context.Context, req domain.CreateEmployee) (models.Employee, error) {
	var employee models.Employee

	err := s.transactor.WithTransaction(ctx, func(ctx context.Context) error {
		_, err := s.positionRepo.CreatePosition(ctx, domain.CreatePosition{
			ID:     req.PositionID,
			Name:   req.PositionName,
			Salary: req.Salary,
		})
		if err != nil {
			return fmt.Errorf("create position: %w", err)
		}

		employee, err = s.employeeRepo.CreateEmployee(ctx, req)
		if err != nil {
			return fmt.Errorf("create employee: %w", err)
		}

		return nil
	})
	if err != nil {
		return models.Employee{}, fmt.Errorf("create employee with transaction: %w", err)
	}

	employeeBytes, err := json.Marshal(employee)

	if err != nil {
		return models.Employee{}, fmt.Errorf("failed to marshal employee: %w", err)
	}

	err = s.eventNotifier.SendMessage(ctx, []byte(employee.ID.String()), employeeBytes)

	if err != nil {
		s.log.Errorf("failed to send message to event notifier: %s", err)
	}

	return employee, nil
}

func (s *EmployeeService) GetEmployee(ctx context.Context, id uuid.UUID) (models.Employee, error) {
	cachedEmployee, err := s.cache.GetEmployee(ctx, id.String())

	if err != nil {
		s.log.Errorf("get employee from cache: %s", err)
	}

	if cachedEmployee != nil {
		s.log.Debugf("returned from cache: %s", cachedEmployee)

		return *cachedEmployee, nil
	}

	employee, err := s.employeeRepo.GetEmployee(ctx, id)
	if err != nil {
		return models.Employee{}, fmt.Errorf("get employee: %w", err)
	}

	if err = s.cache.SetEmployee(ctx, id.String(), &employee); err != nil {
		s.log.Errorf("set employee to cache: %s", err)
	}

	return employee, nil
}

func (s *EmployeeService) GetEmployeeList(ctx context.Context, cursor string) (models.EmployeeList, error) {
	employeeList, err := s.employeeRepo.GetEmployeeList(ctx, cursor)
	if err != nil {
		return models.EmployeeList{}, fmt.Errorf("get employee list: %w", err)
	}

	return employeeList, nil
}

func (s *EmployeeService) UpdateEmployee(ctx context.Context, req domain.UpdateEmployee) (models.Employee, error) {
	employee, err := s.employeeRepo.UpdateEmployee(ctx, req)
	if err != nil {
		return models.Employee{}, fmt.Errorf("update employee: %w", err)
	}

	if err = s.cache.DeleteEmployee(ctx, employee.ID.String()); err != nil {
		s.log.Errorf("delete employee from cache: %s", err)
	}

	return employee, nil
}

func (s *EmployeeService) DeleteEmployee(ctx context.Context, id uuid.UUID) error {
	employee, err := s.employeeRepo.GetEmployee(ctx, id)
	if err != nil {
		return fmt.Errorf("get employee: %w", err)
	}

	err = s.employeeRepo.DeleteEmployee(ctx, employee.ID)
	if err != nil {
		return fmt.Errorf("delete employee: %w", err)
	}

	if err := s.cache.DeleteEmployee(ctx, employee.ID.String()); err != nil {
		s.log.Errorf("delete employee from cache: %s", err)
	}

	return nil
}
