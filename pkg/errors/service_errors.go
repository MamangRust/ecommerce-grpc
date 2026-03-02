package errors

import "fmt"

func NewErrorResponse(message string, code int) error {
	return fmt.Errorf("%s", message)
}
