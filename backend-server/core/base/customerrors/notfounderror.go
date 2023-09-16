package customerrors

import (
	"fmt"
	"net/http"
)

type NotFoundError struct {
	Code    int
	Message string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("[Code: %d] %s", e.Code, e.Message)
}

func NewNotFoundError(message string) *NotFoundError {
	return &NotFoundError{Code: http.StatusNotFound, Message: message}
}
