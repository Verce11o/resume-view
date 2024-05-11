package customerrors

import (
	"errors"
)

var (
	ErrPositionNotFound = errors.New("position not found")
	ErrEmployeeNotFound = errors.New("employee not found")

	ErrEmployeeNotCached = errors.New("employee not cached")
	ErrPositionNotCached = errors.New("position not cached")
)

var ErrDuplicateID = errors.New("duplicate id")
var ErrInvalidCursor = errors.New("invalid cursor")
