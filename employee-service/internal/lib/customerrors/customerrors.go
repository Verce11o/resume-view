package customerrors

import (
	"errors"
)

var ErrPositionNotFound = errors.New("position not found")
var ErrPositionNotUpdated = errors.New("position not updated")

var ErrEmployeeNotFound = errors.New("employee not found")
var ErrEmployeeNotUpdated = errors.New("employee not updated")
