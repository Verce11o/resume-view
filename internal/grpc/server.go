package grpc

import (
	"context"
	"github.com/Verce11o/resume-view/internal/models"
	"github.com/Verce11o/resume-view/lib/grpc_errors"
	pb "github.com/Verce11o/resume-view/protos/gen/go"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type ViewService interface {
	CreateView(ctx context.Context, resumeID, companyID string) (string, error)
	GetResumeViews(ctx context.Context, cursor, resumeID string) (models.ViewList, error)
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
		return nil, status.Errorf(grpc_errors.ParseGRPCErrStatusCode(err), "viewHandler.CreateView: %v", err)
	}

	return &pb.CreateViewResponse{ViewId: viewID}, nil

}

func (s *Server) GetResumeViews(ctx context.Context, request *pb.GetResumeViewsRequest) (*pb.GetResumeViewsResponse, error) {
	ctx, span := s.tracer.Start(ctx, "viewHandler.GetResumeViews")
	defer span.End()

	viewList, err := s.service.GetResumeViews(ctx, request.GetCursor(), request.GetResumeId())
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Errorf(grpc_errors.ParseGRPCErrStatusCode(err), "viewHandler.GetResumeViews: %v", err)
	}

	return viewList.ToProto(), nil
}
