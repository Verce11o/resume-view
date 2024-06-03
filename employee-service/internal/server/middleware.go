package server

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	correlationIDHeader = "X-Correlation-ID"
)

type key int

const (
	keyCorrelationID key = iota
)

func (s *HTTP) LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.log.Debugf("request: %s %s",
			r.Method, r.URL.Path)

		next.ServeHTTP(w, r)
	})
}

func (s *HTTP) CorrelationIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		correlationID := r.Header.Get(correlationIDHeader)
		if correlationID == "" {
			correlationID = uuid.New().String()
		}

		ctx := context.WithValue(r.Context(), keyCorrelationID, correlationID)

		r = r.WithContext(ctx)

		r.Header.Set(correlationIDHeader, correlationID)

		next.ServeHTTP(w, r)
	})
}

func (s *HTTP) TracerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		next.ServeHTTP(w, r)

		duration := time.Since(startTime)

		correlationID, exists := r.Context().Value(keyCorrelationID).(string)
		if !exists {
			correlationID = "unknown"
		}

		s.log.Debugf("correlation ID: %s, request: %s %s, duration: %s",
			correlationID, r.Method, r.URL.Path, duration)
	})
}

func (s *HTTP) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")

		if header == "" {
			http.Error(w, "empty authorization header", http.StatusUnauthorized)

			return
		}

		headerParts := strings.Split(header, " ")

		if len(headerParts) != 2 {
			http.Error(w, "invalid authorization header", http.StatusUnauthorized)

			return
		}

		_, err := s.authenticator.ParseToken(headerParts[1])

		if err != nil {
			http.Error(w, "invalid authorization header", http.StatusUnauthorized)

			return
		}

		next.ServeHTTP(w, r)
	}
}

func (s *HTTP) ContentJSONMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
