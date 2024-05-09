package domain

import "github.com/google/uuid"

type CreatePosition struct {
	ID     uuid.UUID
	Name   string
	Salary int
}

type UpdatePosition struct {
	ID     uuid.UUID
	Name   string
	Salary int
}
