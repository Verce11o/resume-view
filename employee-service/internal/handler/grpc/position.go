package grpc

import (
	"context"

	"github.com/Verce11o/resume-view/employee-service/internal/domain"
	"github.com/Verce11o/resume-view/employee-service/internal/service"
	pb "github.com/Verce11o/resume-view/protos/gen/go"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PositionHandler struct {
	log             *zap.SugaredLogger
	positionService service.Position
	pb.UnimplementedPositionServiceServer
}

func RegisterPosition(server *grpc.Server, log *zap.SugaredLogger, service service.Position) {
	pb.RegisterPositionServiceServer(server, &PositionHandler{log: log, positionService: service})
}

func (h *PositionHandler) CreatePosition(ctx context.Context, input *pb.CreatePositionRequest) (*pb.Position, error) {
	position, err := h.positionService.CreatePosition(ctx, domain.CreatePosition{
		ID:     uuid.New(),
		Name:   input.GetName(),
		Salary: int(input.GetSalary()),
	})

	if err != nil {
		h.log.Errorf("failed to create position: %s", err.Error())

		return nil, status.Errorf(codes.Internal, "failed to create position: %s", err.Error())
	}

	return position.ToProto(), nil
}

func (h *PositionHandler) GetPosition(ctx context.Context, input *pb.GetPositionRequest) (*pb.Position, error) {
	positionID, err := uuid.Parse(input.GetPositionId())
	if err != nil {
		h.log.Errorf("invalid position id: %s", input.GetPositionId())

		return nil, status.Errorf(codes.InvalidArgument, "invalid position id: %s", input.GetPositionId())
	}

	position, err := h.positionService.GetPosition(ctx, positionID)
	if err != nil {
		h.log.Errorf("failed to get position: %s", err.Error())

		return nil, status.Errorf(codes.Internal, "failed to get position: %s", err.Error())
	}

	return position.ToProto(), nil
}

func (h *PositionHandler) GetPositionList(ctx context.Context, input *pb.GetPositionListRequest) (
	*pb.GetPositionListResponse, error) {
	positionList, err := h.positionService.GetPositionList(ctx, input.GetCursor())
	if err != nil {
		h.log.Errorf("failed to get position list: %s", err.Error())

		return nil, status.Errorf(codes.Internal, "failed to get position list: %s", err.Error())
	}

	return positionList.ToProto(), nil
}

func (h *PositionHandler) UpdatePosition(ctx context.Context, input *pb.UpdatePositionRequest) (*pb.Position, error) {
	positionID, err := uuid.Parse(input.GetId())
	if err != nil {
		h.log.Errorf("invalid position id: %s", input.GetId())

		return nil, status.Errorf(codes.InvalidArgument, "invalid position id: %s", input.GetId())
	}

	position, err := h.positionService.UpdatePosition(ctx, domain.UpdatePosition{
		ID:     positionID,
		Name:   input.GetName(),
		Salary: int(input.GetSalary()),
	})

	if err != nil {
		h.log.Errorf("failed to update position: %s", err.Error())

		return nil, status.Errorf(codes.Internal, "failed to update position: %s", err.Error())
	}

	return position.ToProto(), nil
}

func (h *PositionHandler) DeletePosition(ctx context.Context,
	input *pb.DeletePositionRequest) (*pb.DeletePositionResponse, error) {
	positionID, err := uuid.Parse(input.GetPositionId())
	if err != nil {
		h.log.Errorf("invalid position id: %s", input.GetPositionId())

		return nil, status.Errorf(codes.InvalidArgument, "invalid position id: %s", input.GetPositionId())
	}

	err = h.positionService.DeletePosition(ctx, positionID)
	if err != nil {
		h.log.Errorf("failed to delete position: %s", err.Error())

		return nil, status.Errorf(codes.Internal, "failed to delete position: %s", err.Error())
	}

	return &pb.DeletePositionResponse{}, nil
}
