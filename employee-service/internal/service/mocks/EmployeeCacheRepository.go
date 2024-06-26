// Code generated by mockery v2.42.2. DO NOT EDIT.

package mocks

import (
	context "context"

	models "github.com/Verce11o/resume-view/employee-service/internal/models"
	mock "github.com/stretchr/testify/mock"
)

// EmployeeCacheRepository is an autogenerated mock type for the EmployeeCacheRepository type
type EmployeeCacheRepository struct {
	mock.Mock
}

// DeleteEmployee provides a mock function with given fields: ctx, employeeID
func (_m *EmployeeCacheRepository) DeleteEmployee(ctx context.Context, employeeID string) error {
	ret := _m.Called(ctx, employeeID)

	if len(ret) == 0 {
		panic("no return value specified for DeleteEmployee")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, employeeID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetEmployee provides a mock function with given fields: ctx, key
func (_m *EmployeeCacheRepository) GetEmployee(ctx context.Context, key string) (*models.Employee, error) {
	ret := _m.Called(ctx, key)

	if len(ret) == 0 {
		panic("no return value specified for GetEmployee")
	}

	var r0 *models.Employee
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*models.Employee, error)); ok {
		return rf(ctx, key)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *models.Employee); ok {
		r0 = rf(ctx, key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Employee)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetEmployee provides a mock function with given fields: ctx, employeeID, employee
func (_m *EmployeeCacheRepository) SetEmployee(ctx context.Context, employeeID string, employee *models.Employee) error {
	ret := _m.Called(ctx, employeeID, employee)

	if len(ret) == 0 {
		panic("no return value specified for SetEmployee")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, *models.Employee) error); ok {
		r0 = rf(ctx, employeeID, employee)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewEmployeeCacheRepository creates a new instance of EmployeeCacheRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewEmployeeCacheRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *EmployeeCacheRepository {
	mock := &EmployeeCacheRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
