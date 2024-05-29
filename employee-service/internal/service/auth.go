package service

import (
	"context"
	"fmt"

	"github.com/Verce11o/resume-view/employee-service/internal/lib/auth"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AuthService struct {
	log           *zap.SugaredLogger
	employeeRepo  EmployeeRepository
	authenticator *auth.Authenticator
}

func NewAuthService(log *zap.SugaredLogger, employeeRepo EmployeeRepository,
	authenticator *auth.Authenticator) *AuthService {
	return &AuthService{log: log, employeeRepo: employeeRepo, authenticator: authenticator}
}

func (a *AuthService) SignIn(ctx context.Context, employeeID uuid.UUID) (string, error) {
	_, err := a.employeeRepo.GetEmployee(ctx, employeeID)

	if err != nil {
		return "", fmt.Errorf("get employee: %w", err)
	}

	token, err := a.authenticator.GenerateToken(employeeID)

	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}

	return token, nil
}
