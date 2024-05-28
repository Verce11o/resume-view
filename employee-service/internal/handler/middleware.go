package handler

import (
	"context"
	"errors"
	"strings"

	"github.com/getkin/kin-openapi/openapi3filter"
)

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
