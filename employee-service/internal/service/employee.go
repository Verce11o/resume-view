package service

import (
	"context"
	"fmt"
	"github.com/Verce11o/resume-view/employee-service/api"
	"github.com/Verce11o/resume-view/employee-service/internal/models"
	"go.uber.org/zap"
)

type EmployeeRepository interface {
	CreateEmployee(ctx context.Context, request api.CreateEmployee) (models.Employee, error)
	GetEmployee(ctx context.Context, id string) (models.Employee, error)
	GetEmployeeList(ctx context.Context, cursor string) (models.EmployeeList, error)
	UpdateEmployee(ctx context.Context, id string, request api.UpdateEmployee) (models.Employee, error)
	DeleteEmployee(ctx context.Context, id string) error
}

type EmployeeService struct {
	log  *zap.SugaredLogger
	repo EmployeeRepository
}

func NewEmployeeService(log *zap.SugaredLogger, repo EmployeeRepository) *EmployeeService {
	return &EmployeeService{log: log, repo: repo}
}

func (s *EmployeeService) CreateEmployee(ctx context.Context, request api.CreateEmployee) (models.Employee, error) {
	employee, err := s.repo.CreateEmployee(ctx, request)
	if err != nil {
		return models.Employee{}, fmt.Errorf("create employee: %w", err)
	}
	return employee, nil
}

func (s *EmployeeService) GetEmployee(ctx context.Context, id string) (models.Employee, error) {
	employee, err := s.repo.GetEmployee(ctx, id)
	if err != nil {
		return models.Employee{}, fmt.Errorf("get employee: %w", err)
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

func (s *EmployeeService) UpdateEmployee(ctx context.Context, id string, request api.UpdateEmployee) (models.Employee, error) {
	employee, err := s.repo.UpdateEmployee(ctx, id, request)
	if err != nil {
		return models.Employee{}, fmt.Errorf("update employee: %w", err)
	}
	return employee, nil
}

func (s *EmployeeService) DeleteEmployee(ctx context.Context, id string) error {
	err := s.repo.DeleteEmployee(ctx, id)
	if err != nil {
		return fmt.Errorf("delete employee: %w", err)
	}
	return nil
}
