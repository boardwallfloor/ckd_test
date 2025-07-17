package util

import "fmt"

type ServiceError struct {
	Message     string
	ServiceName string
	ErrorMsg    error
}

func (e ServiceError) Error() string {
	return fmt.Sprintf("%s service error. Message: %s, Error Message:`%s`", e.ServiceName, e.Message, e.ErrorMsg)
}
