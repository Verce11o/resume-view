package grpc

import (
	"context"

	pb "github.com/Verce11o/resume-view/protos/gen/go"
	"github.com/Verce11o/resume-view/resume-view/internal/lib/customerrors"
	"github.com/Verce11o/resume-view/resume-view/internal/models"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type ViewService interface {
	CreateView(ctx context.Context, resumeID, companyID string) (uuid.UUID, error)
	ListResumeView(ctx context.Context, cursor, resumeID string) (models.ViewList, error)
}

type Server struct {
	log     *zap.SugaredLogger
	service ViewService
	tracer  trace.Tracer
	pb.UnimplementedViewServiceServer
}

func Register(log *zap.SugaredLogger, service ViewService, server *grpc.Server, tracer trace.Tracer) {
	pb.RegisterViewServiceServer(server, &Server{log: log, service: service, tracer: tracer})
}

func (s *Server) CreateView(ctx context.Context, request *pb.CreateViewRequest) (*pb.CreateViewResponse, error) {
	ctx, span := s.tracer.Start(ctx, "viewHandler.CreateView")
	defer span.End()

	viewID, err := s.service.CreateView(ctx, request.GetResumeId(), request.GetCompanyId())
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return nil, status.Errorf(customerrors.ParseGRPCErrStatusCode(err), "viewHandler.CreateView: %v", err)
	}

	return &pb.CreateViewResponse{ViewId: viewID.String()}, nil
}

func (s *Server) GetResumeViews(ctx context.Context,
	request *pb.GetResumeViewsRequest) (*pb.GetResumeViewsResponse, error) {
	ctx, span := s.tracer.Start(ctx, "viewHandler.GetResumeViews")
	defer span.End()

	viewList, err := s.service.ListResumeView(ctx, request.GetCursor(), request.GetResumeId())
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return nil, status.Errorf(customerrors.ParseGRPCErrStatusCode(err), "viewHandler.ListResumeView: %v", err)
	}

	return viewList.ToProto(), nil
}
