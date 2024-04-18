package models

import (
	"github.com/google/uuid"
	pb "github.com/students-apply/protos/gen/go/view"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type View struct {
	ID        uuid.UUID `json:"id" db:"id"`
	ResumeID  string    `json:"resume_id" db:"resume_id"`
	CompanyID uuid.UUID `json:"company_id" db:"company_id"`
	ViewedAt  time.Time `json:"viewed_at" db:"viewed_at"`
}

func (v *View) ToProto() *pb.View {
	return &pb.View{
		ViewId:    v.ID.String(),
		ResumeId:  v.ResumeID,
		CompanyId: v.CompanyID.String(),
		ViewedAt:  timestamppb.New(v.ViewedAt),
	}
}

type ViewList struct {
	Cursor string `json:"cursor"`
	Views  []View `json:"views"`
	Total  int    `json:"total"`
}

func (v *ViewList) ToProto() *pb.GetResumeViewsResponse {
	views := make([]*pb.View, 0, len(v.Views))
	for _, val := range v.Views {
		views = append(views, val.ToProto())
	}
	return &pb.GetResumeViewsResponse{
		Views:  views,
		Cursor: v.Cursor,
		Total:  int32(v.Total),
	}
}
