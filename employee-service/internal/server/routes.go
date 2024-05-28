package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Verce11o/resume-view/employee-service/api"
	"github.com/Verce11o/resume-view/employee-service/internal/config"
	"github.com/Verce11o/resume-view/employee-service/internal/handler"
	"github.com/Verce11o/resume-view/employee-service/internal/lib/auth"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/gin-gonic/gin"
	middleware "github.com/oapi-codegen/gin-middleware"
	"go.uber.org/zap"
)

type Server struct {
	log             *zap.SugaredLogger
	employeeService handler.EmployeeService
	positionService handler.PositionService
	authenticator   *auth.Authenticator
	cfg             config.Config
	httpServer      *http.Server
}

func NewServer(log *zap.SugaredLogger, employeeService handler.EmployeeService, positionService handler.PositionService,
	cfg config.Config, authenticator *auth.Authenticator) *Server {
	return &Server{log: log, employeeService: employeeService, positionService: positionService,
		authenticator: authenticator, cfg: cfg}
}

func (s *Server) Run(handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:         s.cfg.HTTPServer.Port,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	s.log.Infof("HTTPServer running on: %v", s.cfg.HTTPServer.Port)

	if err := s.httpServer.ListenAndServe(); err != nil {
		return fmt.Errorf("failed to start HTTP server: %w", err)
	}

	return nil
}

func (s *Server) InitRoutes() *gin.Engine {
	router := gin.New()

	apiGroup := router.Group("/api/v1")

	spec, _ := api.GetSwagger()
	handlers := handler.NewHandler(s.log, s.positionService, s.employeeService, s.authenticator)

	validator := middleware.OapiRequestValidatorWithOptions(spec,
		&middleware.Options{
			ErrorHandler: func(c *gin.Context, message string, statusCode int) {
				c.JSON(http.StatusUnauthorized, gin.H{
					"message": message,
				})
			},
			Options: openapi3filter.Options{
				AuthenticationFunc: handlers.AuthMiddleware,
			},
		},
	)

	apiGroup.Use(validator)

	api.RegisterHandlers(apiGroup, handlers)

	return router
}

func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown http server: %w", err)
	}

	return nil
}
