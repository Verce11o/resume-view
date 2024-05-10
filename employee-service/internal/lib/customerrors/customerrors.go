package customerrors

import (
	"errors"
)

var ErrPositionNotFound = errors.New("position not found")

var ErrEmployeeNotFound = errors.New("employee not found")

var ErrDuplicateID = errors.New("duplicate id")

var ErrInvalidCursor = errors.New("invalid cursor")
