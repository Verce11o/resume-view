package server

import (
	"context"
	"github.com/Verce11o/resume-view/employee-service/api"
	"github.com/Verce11o/resume-view/employee-service/internal/config"
	"github.com/Verce11o/resume-view/employee-service/internal/handler"
	"github.com/Verce11o/resume-view/employee-service/internal/repository/postgres"
	"github.com/Verce11o/resume-view/employee-service/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Server struct {
	log        *zap.SugaredLogger
	db         *pgxpool.Pool
	cfg        config.Config
	httpServer *http.Server
}

func NewServer(log *zap.SugaredLogger, db *pgxpool.Pool, cfg config.Config) *Server {
	return &Server{log: log, db: db, cfg: cfg}
}

func (s *Server) Run(handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:         s.cfg.Server.Port,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	s.log.Infof("Server running on: %v", s.cfg.Server.Port)
	return s.httpServer.ListenAndServe()
}

func (s *Server) InitRoutes() *gin.Engine {

	router := gin.New()

	positionRepo := postgres.NewPositionRepository(s.db)

	positionService := service.NewPositionService(s.log, positionRepo)

	apiGroup := router.Group("/api/v1")

	api.RegisterHandlers(apiGroup, handler.NewHandler(s.log, positionService, nil))

	return router
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
