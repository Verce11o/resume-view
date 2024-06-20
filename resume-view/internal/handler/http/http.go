package http

import (
	"errors"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

type Server struct {
	log        *zap.SugaredLogger
	httpServer *http.Server
	port       string
}

func NewServer(log *zap.SugaredLogger, port string) *Server {
	return &Server{log: log, port: port}
}

func (s *Server) Run() error {
	s.httpServer = &http.Server{
		Addr:         s.port,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	http.Handle("/metrics", promhttp.Handler())

	if err := s.httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		s.log.Fatalf("http server error: %v", err)
	}

	return nil
}
