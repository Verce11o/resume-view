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

type EmployeeHandler struct {
	log             *zap.SugaredLogger
	employeeService service.Employee
	pb.UnimplementedEmployeeServiceServer
}

func RegisterEmployee(server *grpc.Server, log *zap.SugaredLogger, service service.Employee) {
	pb.RegisterEmployeeServiceServer(server, &EmployeeHandler{log: log, employeeService: service})
}

func (h *EmployeeHandler) CreateEmployee(ctx context.Context, input *pb.CreateEmployeeRequest) (*pb.Employee, error) {
	employee, err := h.employeeService.CreateEmployee(ctx, domain.CreateEmployee{
		EmployeeID:   uuid.New(),
		PositionID:   uuid.New(),
		FirstName:    input.GetFirstName(),
		LastName:     input.GetLastName(),
		PositionName: input.GetPositionName(),
		Salary:       int(input.GetSalary()),
	})

	if err != nil {
		h.log.Errorf("failed to create employee: %s", err.Error())

		return nil, status.Errorf(codes.Internal, "failed to create employee: %s", err.Error())
	}

	return employee.ToProto(), nil
}

func (h *EmployeeHandler) GetEmployee(ctx context.Context, input *pb.GetEmployeeRequest) (*pb.Employee, error) {
	employeeID, err := uuid.Parse(input.GetEmployeeId())
	if err != nil {
		h.log.Errorf("invalid employee id: %s", input.GetEmployeeId())

		return nil, status.Errorf(codes.InvalidArgument, "invalid employee id: %s", input.GetEmployeeId())
	}

	employee, err := h.employeeService.GetEmployee(ctx, employeeID)
	if err != nil {
		h.log.Errorf("failed to get employee: %s", err.Error())

		return nil, status.Errorf(codes.Internal, "failed to get employee: %s", err.Error())
	}

	return employee.ToProto(), nil
}

func (h *EmployeeHandler) GetEmployeeList(ctx context.Context, input *pb.GetEmployeeListRequest) (
	*pb.GetEmployeeListResponse, error) {
	employeeList, err := h.employeeService.GetEmployeeList(ctx, input.GetCursor())
	if err != nil {
		h.log.Errorf("failed to get employee list: %s", err.Error())

		return nil, status.Errorf(codes.Internal, "failed to get employee list: %s", err.Error())
	}

	return employeeList.ToProto(), nil
}

func (h *EmployeeHandler) UpdateEmployee(ctx context.Context, input *pb.UpdateEmployeeRequest) (*pb.Employee, error) {
	employeeID, err := uuid.Parse(input.GetEmployeeId())
	if err != nil {
		h.log.Errorf("invalid employee id: %s", input.GetEmployeeId())

		return nil, status.Errorf(codes.InvalidArgument, "invalid employee id: %s", input.GetEmployeeId())
	}

	positionID, err := uuid.Parse(input.GetPositionId())
	if err != nil {
		h.log.Errorf("invalid position id: %s", input.GetPositionId())

		return nil, status.Errorf(codes.InvalidArgument, "invalid position id: %s", input.GetPositionId())
	}

	employee, err := h.employeeService.UpdateEmployee(ctx, domain.UpdateEmployee{
		EmployeeID: employeeID,
		PositionID: positionID,
		FirstName:  input.GetFirstName(),
		LastName:   input.GetLastName(),
		Salary:     int(input.GetSalary()),
	})

	if err != nil {
		h.log.Errorf("failed to update employee: %s", err.Error())

		return nil, status.Errorf(codes.Internal, "failed to update employee: %s", err.Error())
	}

	return employee.ToProto(), nil
}

func (h *EmployeeHandler) DeleteEmployee(ctx context.Context,
	input *pb.DeleteEmployeeRequest) (*pb.DeleteEmployeeResponse, error) {
	employeeID, err := uuid.Parse(input.GetEmployeeId())
	if err != nil {
		h.log.Errorf("invalid employee id: %s", input.GetEmployeeId())

		return nil, status.Errorf(codes.InvalidArgument, "invalid employee id: %s", input.GetEmployeeId())
	}

	err = h.employeeService.DeleteEmployee(ctx, employeeID)
	if err != nil {
		h.log.Errorf("failed to delete employee: %s", err.Error())

		return nil, status.Errorf(codes.Internal, "failed to delete employee: %s", err.Error())
	}

	return &pb.DeleteEmployeeResponse{}, nil
}
