package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Verce11o/resume-view/employee-service/api"
	"github.com/Verce11o/resume-view/employee-service/internal/config"
	http2 "github.com/Verce11o/resume-view/employee-service/internal/handler/http"
	"github.com/Verce11o/resume-view/employee-service/internal/lib/auth"
	"github.com/Verce11o/resume-view/employee-service/internal/service"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/gin-gonic/gin"
	middleware "github.com/oapi-codegen/gin-middleware"
	"go.uber.org/zap"
)

type HTTP struct {
	log             *zap.SugaredLogger
	employeeService service.Employee
	positionService service.Position
	authenticator   *auth.Authenticator
	cfg             config.Config
	httpServer      *http.Server
}

func NewHTTP(log *zap.SugaredLogger, employeeService service.Employee, positionService service.Position,
	authenticator *auth.Authenticator, cfg config.Config) *HTTP {
	return &HTTP{log: log, employeeService: employeeService, positionService: positionService,
		authenticator: authenticator, cfg: cfg}
}

func (s *HTTP) Run(handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:         s.cfg.HTTPServer.Port,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		if err := s.httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			s.log.Fatalf("HTTPServer error: %v", err)
		}
	}()

	return nil
}

func (s *HTTP) InitRoutes() *gin.Engine {
	router := gin.New()

	apiGroup := router.Group("/api/v1")

	spec, _ := api.GetSwagger()
	handlers := http2.NewHandler(s.log, s.positionService, s.employeeService, s.authenticator)

	validator := middleware.OapiRequestValidatorWithOptions(spec,
		&middleware.Options{
			ErrorHandler: func(c *gin.Context, message string, _ int) {
				c.JSON(http.StatusUnauthorized, gin.H{
					"message": message,
				})
			},
			Options: openapi3filter.Options{
				AuthenticationFunc: handlers.AuthMiddleware,
			},
		},
	)

	apiGroup.Use(handlers.LogMiddleware)
	apiGroup.Use(handlers.CorrelationIDMiddleware)
	apiGroup.Use(handlers.TracerMiddleware)
	apiGroup.Use(validator)

	api.RegisterHandlers(apiGroup, handlers)

	return router
}

func (s *HTTP) Shutdown(ctx context.Context) error {
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown http server: %w", err)
	}

	return nil
}
