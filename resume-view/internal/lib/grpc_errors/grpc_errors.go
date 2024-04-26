package grpc_errors

import (
	"context"
	"errors"
	"google.golang.org/grpc/codes"
)

var (
	ErrNotFound      = errors.New("view not found")
	ErrInvalidCursor = errors.New("invalid cursor")
)

func ParseGRPCErrStatusCode(err error) codes.Code {
	switch {
	case errors.Is(err, context.Canceled):
		return codes.Canceled
	case errors.Is(err, context.DeadlineExceeded):
		return codes.DeadlineExceeded
	case errors.Is(err, ErrNotFound):
		return codes.NotFound
	case errors.Is(err, ErrInvalidCursor):
		return codes.InvalidArgument
	}
	return codes.Internal
}
