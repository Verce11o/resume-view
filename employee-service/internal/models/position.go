package models

import (
	"github.com/google/uuid"
	"time"
)

type Position struct {
	ID        uuid.UUID `json:"id" db:"id" bson:"_id,omitempty"`
	Name      string    `json:"name" db:"name" bson:"name,omitempty"`
	Salary    int       `json:"salary" db:"salary" bson:"salary,omitempty"`
	CreatedAt time.Time `json:"created_at" db:"created_at" bson:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" bson:"updated_at,omitempty"`
}

type PositionList struct {
	Cursor    string     `json:"cursor"`
	Positions []Position `json:"positions"`
}
