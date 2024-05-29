package models

import (
	"time"

	pb "github.com/Verce11o/resume-view/protos/gen/go"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Position struct {
	ID        uuid.UUID `json:"id" db:"id" bson:"_id,omitempty"`
	Name      string    `json:"name" db:"name" bson:"name,omitempty"`
	Salary    int       `json:"salary" db:"salary" bson:"salary,omitempty"`
	CreatedAt time.Time `json:"created_at" db:"created_at" bson:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" bson:"updated_at,omitempty"`
}

func (p *Position) ToProto() *pb.Position {
	return &pb.Position{
		Id:        p.ID.String(),
		Name:      p.Name,
		Salary:    int32(p.Salary),
		CreatedAt: timestamppb.New(p.CreatedAt),
		UpdatedAt: timestamppb.New(p.UpdatedAt),
	}
}

type PositionList struct {
	Cursor    string     `json:"cursor"`
	Positions []Position `json:"positions"`
}

func (p *PositionList) ToProto() *pb.GetPositionListResponse {
	positions := make([]*pb.Position, 0, len(p.Positions))
	for _, val := range p.Positions {
		positions = append(positions, val.ToProto())
	}
	return &pb.GetPositionListResponse{
		Cursor:    p.Cursor,
		Positions: positions,
	}
}
