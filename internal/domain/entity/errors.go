package entity

import (
	"fmt"
)

// ErrInvalidInput represents an error due to invalid input data
type ErrInvalidInput struct {
	Field   string
	Message string
}

// Error returns the error message
func (e ErrInvalidInput) Error() string {
	return fmt.Sprintf("invalid input for %s: %s", e.Field, e.Message)
}
