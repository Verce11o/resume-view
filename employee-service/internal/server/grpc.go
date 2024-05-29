package server

import (
	"fmt"
	"net"

	"github.com/Verce11o/resume-view/employee-service/internal/config"
	employeeGrpc "github.com/Verce11o/resume-view/employee-service/internal/handler/grpc"
	"github.com/Verce11o/resume-view/employee-service/internal/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type GRPC struct {
	log             *zap.SugaredLogger
	employeeService service.Employee
	positionService service.Position
	cfg             config.Config
	server          *grpc.Server
}

func NewGRPC(log *zap.SugaredLogger, employeeService service.Employee, positionService service.Position, cfg config.Config) *GRPC {
	srv := grpc.NewServer()
	return &GRPC{log: log, employeeService: employeeService, positionService: positionService, cfg: cfg, server: srv}
}

func (g *GRPC) Run() error {
	employeeGrpc.RegisterEmployee(g.server, g.log, g.employeeService)
	employeeGrpc.RegisterPosition(g.server, g.log, g.positionService)

	l, err := net.Listen("tcp", g.cfg.GRPCServer.Port)

	if err != nil {
		return fmt.Errorf("failed to listen to tcp: %w", err)
	}

	if err := g.server.Serve(l); err != nil {
		return fmt.Errorf("failed to serve grpc: %w", err)
	}

	return nil
}

func (g *GRPC) Shutdown() {
	g.server.GracefulStop()
}
