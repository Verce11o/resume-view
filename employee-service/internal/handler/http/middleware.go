package http

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	correlationIDHeader = "X-Correlation-ID"
	correlationIDCtx    = "correlation-id"
)

func (h *Handler) LogMiddleware(c *gin.Context) {
	h.log.Debugf("request: %s %s, status: %d",
		c.Request.Method, c.Request.URL.Path, c.Writer.Status())

	c.Next()
}

func (h *Handler) CorrelationIDMiddleware(c *gin.Context) {
	correlationID := c.Request.Header.Get(correlationIDHeader)
	if correlationID == "" {
		correlationID = uuid.New().String()
	}

	c.Set(correlationIDCtx, correlationID)
	c.Writer.Header().Set(correlationIDHeader, correlationID)
	c.Next()
}

func (h *Handler) TracerMiddleware(c *gin.Context) {
	startTime := time.Now()

	c.Next()

	duration := time.Since(startTime)

	correlationID, exists := c.Get(correlationIDCtx)
	if !exists {
		correlationID = "unknown"
	}

	h.log.Debugf("correlation ID: %s, request: %s %s, duration: %s",
		correlationID, c.Request.Method, c.Request.URL.Path, duration)
}

func (h *Handler) AuthMiddleware(_ context.Context, input *openapi3filter.AuthenticationInput) error {
	req := input.RequestValidationInput.Request
	header := req.Header.Get("Authorization")

	if header == "" {
		return errors.New("authorization header is required")
	}

	headerParts := strings.Split(header, " ")

	if len(headerParts) != 2 {
		return errors.New("invalid authorization header")
	}

	_, err := h.authenticator.ParseToken(headerParts[1])

	if err != nil {
		return errors.New("invalid authorization header")
	}

	return nil
}
