package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/rs/cors"
	"net/http"
	"time"

	"github.com/Verce11o/resume-view/employee-service/internal/config"
	"github.com/Verce11o/resume-view/employee-service/internal/handler"
	chiHandler "github.com/Verce11o/resume-view/employee-service/internal/handler/http/chi"
	"github.com/Verce11o/resume-view/employee-service/internal/handler/http/gorilla"
	"github.com/Verce11o/resume-view/employee-service/internal/lib/auth"
	"github.com/Verce11o/resume-view/employee-service/internal/service"
	"github.com/go-chi/chi"
	gorillaMux "github.com/gorilla/mux"
	"go.uber.org/zap"
)

type HTTP struct {
	log             *zap.SugaredLogger
	employeeService service.Employee
	positionService service.Position
	authService     service.Auth
	authenticator   *auth.Authenticator
	cfg             config.Config
	httpServer      *http.Server
}

func NewHTTP(log *zap.SugaredLogger, employeeService service.Employee, positionService service.Position,
	authService service.Auth, authenticator *auth.Authenticator, cfg config.Config) *HTTP {
	return &HTTP{log: log, employeeService: employeeService, positionService: positionService,
		authService: authService, authenticator: authenticator, cfg: cfg}
}

func (s *HTTP) Run(handler http.Handler) error {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5174"},                                                                // Specifically allow your frontend origin
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions}, // Added OPTIONS for preflight
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum age (in seconds) of the preflight request cache
		ExposedHeaders:   []string{"Link"},
		Debug:            true, // Enable debug mode to help diagnose issues (remove in production)
	})

	handler = c.Handler(handler)

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

func (s *HTTP) Shutdown(ctx context.Context) error {
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown http server: %w", err)
	}

	return nil
}

type CustomRouter interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	MethodFunc(method, path string, handler http.HandlerFunc)
	Use(middleware ...func(http.Handler) http.Handler)
}

func (s *HTTP) InitRoutes() (CustomRouter, error) {
	var (
		router          CustomRouter
		employeeHandler handler.EmployeeHandler
		positionHandler handler.PositionHandler
	)

	switch s.cfg.HTTPServer.Router {
	case "chi":
		router = chi.NewRouter()
		h := chiHandler.New(s.log, s.positionService, s.employeeService, s.authService)
		employeeHandler, positionHandler = h, h

	case "gorilla":
		router = gorilla.NewWrappedRouter(gorillaMux.NewRouter())
		h := gorilla.New(s.log, s.positionService, s.employeeService, s.authService)
		employeeHandler, positionHandler = h, h

	default:
		return nil, fmt.Errorf("invalid router type: %s", s.cfg.HTTPServer.Router)
	}

	router.Use(s.LogMiddleware, s.CorrelationIDMiddleware, s.ContentJSONMiddleware)

	router.MethodFunc(http.MethodPost, "/auth/signin", employeeHandler.SignIn)

	{
		router.MethodFunc(http.MethodGet, "/employee", employeeHandler.GetEmployeeList)
		router.MethodFunc(http.MethodPost, "/employee", s.AuthMiddleware(employeeHandler.CreateEmployee))
		router.MethodFunc(http.MethodGet, "/employee/{id}", employeeHandler.GetEmployeeByID)
		router.MethodFunc(http.MethodPut, "/employee/{id}", s.AuthMiddleware(employeeHandler.UpdateEmployeeByID))
		router.MethodFunc(http.MethodDelete, "/employee/{id}", s.AuthMiddleware(employeeHandler.DeleteEmployeeByID))
	}

	{
		router.MethodFunc(http.MethodGet, "/position", positionHandler.GetPositionList)
		router.MethodFunc(http.MethodPost, "/position", s.AuthMiddleware(positionHandler.CreatePosition))
		router.MethodFunc(http.MethodGet, "/position/{id}", positionHandler.GetPositionByID)
		router.MethodFunc(http.MethodPut, "/position/{id}", s.AuthMiddleware(positionHandler.UpdatePositionByID))
		router.MethodFunc(http.MethodDelete, "/position/{id}", s.AuthMiddleware(positionHandler.DeletePositionByID))
	}

	return router, nil
}
