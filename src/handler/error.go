package handler

import (
	"fmt"
)

type ApiError struct {
	Status int
	Msg    any
	Err    error
}

func (e ApiError) Error() string {
	return fmt.Sprintf("%s", e.Msg)
}

func NewApiError(status int, msg any, err error) ApiError {
	return ApiError{
		Status: status,
		Msg:    msg,
		Err:    err,
	}
}

const (
	EMAIL_ALREADY_EXIST       = "The email address has already been taken."
	INVALID_PAYLOAD           = "Invalid JSON payload."
	INVALID_LOGIN_CREDENTIALS = "invalidLoginCredentials"
	UNAUTHORIZED_REQUEST      = "Unauthorized request"
)
