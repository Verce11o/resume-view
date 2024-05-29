package service

import (
	"context"

	"github.com/Verce11o/resume-view/employee-service/internal/domain"
	"github.com/Verce11o/resume-view/employee-service/internal/models"
	"github.com/google/uuid"
)

//go:generate mockgen -source=service.go -destination=mocks/services.go -package=mocks -mock_names=Employee=MockEmployeeService,Position=MockPositionService
type Employee interface {
	CreateEmployee(ctx context.Context, req domain.CreateEmployee) (models.Employee, error)
	GetEmployee(ctx context.Context, id uuid.UUID) (models.Employee, error)
	GetEmployeeList(ctx context.Context, cursor string) (models.EmployeeList, error)
	UpdateEmployee(ctx context.Context, req domain.UpdateEmployee) (models.Employee, error)
	DeleteEmployee(ctx context.Context, id uuid.UUID) error
	SignIn(ctx context.Context, employeeID uuid.UUID) (string, error)
}

type Position interface {
	CreatePosition(ctx context.Context, req domain.CreatePosition) (models.Position, error)
	GetPosition(ctx context.Context, id uuid.UUID) (models.Position, error)
	GetPositionList(ctx context.Context, cursor string) (models.PositionList, error)
	UpdatePosition(ctx context.Context, req domain.UpdatePosition) (models.Position, error)
	DeletePosition(ctx context.Context, id uuid.UUID) error
}
