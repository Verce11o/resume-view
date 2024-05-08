package service

import (
	"context"
	"fmt"

	"github.com/Verce11o/resume-view/employee-service/api"
	"github.com/Verce11o/resume-view/employee-service/internal/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type EmployeeRepository interface {
	CreateEmployee(ctx context.Context, employeeID uuid.UUID, positionID uuid.UUID, request api.CreateEmployee) (models.Employee, error)
	GetEmployee(ctx context.Context, id uuid.UUID) (models.Employee, error)
	GetEmployeeList(ctx context.Context, cursor string) (models.EmployeeList, error)
	UpdateEmployee(ctx context.Context, id uuid.UUID, request api.UpdateEmployee) (models.Employee, error)
	DeleteEmployee(ctx context.Context, id uuid.UUID) error
}

type EmployeeCacheRepository interface {
	GetEmployee(ctx context.Context, key string) (*models.Employee, error)
	SetEmployee(ctx context.Context, employeeID string, employee *models.Employee) error
	DeleteEmployee(ctx context.Context, employeeID string) error
}

type EmployeeService struct {
	log   *zap.SugaredLogger
	repo  EmployeeRepository
	cache EmployeeCacheRepository
}

func NewEmployeeService(log *zap.SugaredLogger, repo EmployeeRepository, cache EmployeeCacheRepository) *EmployeeService {
	return &EmployeeService{log: log, repo: repo, cache: cache}
}

func (s *EmployeeService) CreateEmployee(ctx context.Context, employeeID uuid.UUID, positionID uuid.UUID, request api.CreateEmployee) (models.Employee, error) {
	employee, err := s.repo.CreateEmployee(ctx, employeeID, positionID, request)
	if err != nil {
		return models.Employee{}, fmt.Errorf("create employee: %w", err)
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

	employee, err := s.repo.GetEmployee(ctx, id)
	if err != nil {
		return models.Employee{}, fmt.Errorf("get employee: %w", err)
	}

	if err = s.cache.SetEmployee(ctx, id.String(), &employee); err != nil {
		s.log.Errorf("set employee to cache: %s", err)
	}

	return employee, nil
}

func (s *EmployeeService) GetEmployeeList(ctx context.Context, cursor string) (models.EmployeeList, error) {
	employeeList, err := s.repo.GetEmployeeList(ctx, cursor)
	if err != nil {
		return models.EmployeeList{}, fmt.Errorf("get employee list: %w", err)
	}
	return employeeList, nil
}

func (s *EmployeeService) UpdateEmployee(ctx context.Context, id uuid.UUID, request api.UpdateEmployee) (models.Employee, error) {
	employee, err := s.repo.UpdateEmployee(ctx, id, request)
	if err != nil {
		return models.Employee{}, fmt.Errorf("update employee: %w", err)
	}

	if err = s.cache.DeleteEmployee(ctx, employee.ID.String()); err != nil {
		s.log.Errorf("delete employee from cache: %s", err)
	}

	return employee, nil
}

func (s *EmployeeService) DeleteEmployee(ctx context.Context, id uuid.UUID) error {
	employee, err := s.repo.GetEmployee(ctx, id)
	if err != nil {
		return fmt.Errorf("get employee: %w", err)
	}

	err = s.repo.DeleteEmployee(ctx, employee.ID)
	if err != nil {
		return fmt.Errorf("delete employee: %w", err)
	}

	if err := s.cache.DeleteEmployee(ctx, employee.ID.String()); err != nil {
		s.log.Errorf("delete employee from cache: %s", err)
	}

	return nil
}
