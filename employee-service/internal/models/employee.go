package models

import (
	"github.com/google/uuid"
	"time"
)

type Employee struct {
	ID         uuid.UUID `json:"id" db:"id" bson:"_id,omitempty"`
	FirstName  string    `json:"first_name" db:"first_name" bson:"first_name,omitempty"`
	LastName   string    `json:"last_name" db:"last_name" bson:"last_name,omitempty"`
	PositionID uuid.UUID `json:"position_id" db:"position_id" bson:"position_id,omitempty"`
	CreatedAt  time.Time `json:"created_at" db:"created_at" bson:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at" bson:"updated_at,omitempty"`
}

type EmployeeList struct {
	Cursor    string     `json:"cursor"`
	Employees []Employee `json:"employees"`
}
