package models

import (
	"github.com/google/uuid"
	"time"
)

type Employee struct {
	ID         uuid.UUID `json:"id" db:"id"`
	FirstName  string    `json:"first_name" db:"first_name"`
	LastName   string    `json:"last_name" db:"last_name"`
	PositionID uuid.UUID `json:"position_id" db:"position_id"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

type EmployeeList struct {
	Cursor    string     `json:"cursor"`
	Employees []Employee `json:"employees"`
}
