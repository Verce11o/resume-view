package models

import (
	"time"

	pb "github.com/Verce11o/resume-view/protos/gen/go"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Employee struct {
	ID         uuid.UUID `json:"id" db:"id" bson:"_id,omitempty"`
	FirstName  string    `json:"first_name" db:"first_name" bson:"first_name,omitempty"`
	LastName   string    `json:"last_name" db:"last_name" bson:"last_name,omitempty"`
	PositionID uuid.UUID `json:"position_id" db:"position_id" bson:"position_id,omitempty"`
	CreatedAt  time.Time `json:"created_at" db:"created_at" bson:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at" bson:"updated_at,omitempty"`
}

func (e *Employee) ToProto() *pb.Employee {
	return &pb.Employee{
		Id:         e.ID.String(),
		FirstName:  e.FirstName,
		LastName:   e.LastName,
		PositionId: e.PositionID.String(),
		CreatedAt:  timestamppb.New(e.CreatedAt),
		UpdatedAt:  timestamppb.New(e.UpdatedAt),
	}
}

type EmployeeList struct {
	Cursor    string     `json:"cursor"`
	Employees []Employee `json:"employees"`
}

func (e *EmployeeList) ToProto() *pb.GetEmployeeListResponse {
	employees := make([]*pb.Employee, 0, len(e.Employees))
	for _, val := range e.Employees {
		employees = append(employees, val.ToProto())
	}

	return &pb.GetEmployeeListResponse{
		Cursor:    e.Cursor,
		Employees: employees,
	}
}
