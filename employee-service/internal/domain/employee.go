package domain

import "github.com/google/uuid"

type CreateEmployee struct {
	EmployeeID   uuid.UUID
	PositionID   uuid.UUID
	FirstName    string
	LastName     string
	PositionName string
	Salary       int
}

type UpdateEmployee struct {
	EmployeeID uuid.UUID
	PositionID uuid.UUID
	FirstName  string
	LastName   string
	Salary     int
}
