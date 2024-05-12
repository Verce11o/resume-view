package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Verce11o/resume-view/employee-service/api"
	"github.com/Verce11o/resume-view/employee-service/internal/config"
	"github.com/Verce11o/resume-view/employee-service/internal/handler"
	"github.com/Verce11o/resume-view/employee-service/internal/repository/mongodb"
	"github.com/Verce11o/resume-view/employee-service/internal/repository/postgres"
	"github.com/Verce11o/resume-view/employee-service/internal/repository/redis"
	"github.com/Verce11o/resume-view/employee-service/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	rdb "github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type Server struct {
	log        *zap.SugaredLogger
	db         *pgxpool.Pool
	mongo      *mongo.Client
	redis      *rdb.Client
	cfg        config.Config
	httpServer *http.Server
}

func NewServer(
	log *zap.SugaredLogger,
	db *pgxpool.Pool,
	mongo *mongo.Client,
	redis *rdb.Client,
	cfg config.Config) *Server {
	return &Server{log: log, db: db, mongo: mongo, redis: redis, cfg: cfg}
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

	var (
		positionRepo service.PositionRepository
		employeeRepo service.EmployeeRepository
		transactor   service.Transactor
	)

	positionRepo = postgres.NewPositionRepository(s.db)
	employeeRepo = postgres.NewEmployeeRepository(s.db)
	transactor = postgres.NewTransactor(s.db)

	if strings.ToLower(s.cfg.MainDatabase) == "mongo" {
		positionRepo = mongodb.NewPositionRepository(s.mongo.Database("employees"))
		employeeRepo = mongodb.NewEmployeeRepository(s.mongo.Database("employees"))
		transactor = mongodb.NewTransactor(s.mongo)
	}

	positionCache := redis.NewPositionCache(s.redis)
	employeeCache := redis.NewEmployeeCache(s.redis)

	positionService := service.NewPositionService(s.log, positionRepo, positionCache)

	employeeService := service.NewEmployeeService(s.log, employeeRepo, positionRepo, employeeCache, transactor)

	apiGroup := router.Group("/api/v1")

	api.RegisterHandlers(apiGroup, handler.NewHandler(s.log, positionService, employeeService))

	return router
}

func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown http server: %w", err)
	}

	return nil
}
