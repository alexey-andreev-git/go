package error

import "fmt"

// Define your custom error types here
type CustomError struct {
	Message string
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("Custom Error: %s", e.Message)
}
