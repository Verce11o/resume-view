package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Verce11o/resume-view/employee-service/internal/models"
	"github.com/redis/go-redis/v9"
)

const (
	employeeTTL = 3600
)

type EmployeeCache struct {
	client *redis.Client
}

func NewEmployeeCache(client *redis.Client) *EmployeeCache {
	return &EmployeeCache{client: client}
}

func (r *EmployeeCache) GetEmployee(ctx context.Context, employeeID string) (*models.Employee, error) {
	employeeBytes, err := r.client.Get(ctx, r.createKey(employeeID)).Bytes()

	if err != nil || errors.Is(err, redis.Nil) {
		return nil, err
	}

	var employee models.Employee

	if err = json.Unmarshal(employeeBytes, &employee); err != nil {
		return nil, err
	}

	return &employee, nil
}

func (r *EmployeeCache) SetEmployee(ctx context.Context, employeeID string, employee *models.Employee) error {
	employeeBytes, err := json.Marshal(employee)

	if err != nil {
		return err
	}

	return r.client.Set(ctx, r.createKey(employeeID), employeeBytes, time.Second*time.Duration(employeeTTL)).Err()
}

func (r *EmployeeCache) DeleteEmployee(ctx context.Context, employeeID string) error {
	return r.client.Del(ctx, r.createKey(employeeID)).Err()
}

func (r *EmployeeCache) createKey(key string) string {
	return fmt.Sprintf("employee:%s", key)
}
