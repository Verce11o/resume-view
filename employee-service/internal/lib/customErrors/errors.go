package customErrors

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

var ErrNotFound = errors.New("not found")

func MapHTTPError(err error) (int, string) {
	switch {
	case errors.Is(err, ErrNotFound):
		return http.StatusNotFound, err.Error()
	}
	return http.StatusInternalServerError, "server error"
}

func WithHTTPError(c *gin.Context, err error) {
	status, message := MapHTTPError(err)
	c.JSON(status, gin.H{
		"message": message,
	})
	return
}
