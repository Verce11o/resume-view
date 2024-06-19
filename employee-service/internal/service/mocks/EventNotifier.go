// Code generated by mockery v2.42.2. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// EventNotifier is an autogenerated mock type for the EventNotifier type
type EventNotifier struct {
	mock.Mock
}

// SendMessage provides a mock function with given fields: ctx, key, value
func (_m *EventNotifier) SendMessage(ctx context.Context, key []byte, value []byte) error {
	ret := _m.Called(ctx, key, value)

	if len(ret) == 0 {
		panic("no return value specified for SendMessage")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []byte, []byte) error); ok {
		r0 = rf(ctx, key, value)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewEventNotifier creates a new instance of EventNotifier. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewEventNotifier(t interface {
	mock.TestingT
	Cleanup(func())
}) *EventNotifier {
	mock := &EventNotifier{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
