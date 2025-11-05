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
	UNAUTHORIZED_REQUEST = "Unauthorized request"
	INVALID_PAYLOAD      = "Invalid JSON payload."
	UPDATE_BAD_REQUEST   = "update data failed"

	EMAIL_ALREADY_EXIST       = "The email address has already been taken."
	INVALID_LOGIN_CREDENTIALS = "invalidLoginCredentials"

	ERROR_RETRIEVING_USER_FURNITURE_LIST  = "Error retrieving user furniture list"
	ERROR_RETRIEVING_FURNITURE_CATEGORIES = "Error retrieving furniture categories"
	FURNITURE_NOT_ACCESSED                = "Furniture not found or access denied"
	ERROR_DELETING_USER_FURNITURE         = "Failed to delete furniture"

	ERROR_RETRIEVING_USER_BLUEPRINTS = "Error retrieving user blueprints"
	INVALID_BLUEPRINT_ID             = "Invalid Blueprint ID"
	BLUEPRINT_NOT_ACCESSED           = "Blueprint not found or access denied"
	ERROR_DELETING_USER_BLUEPRINT    = "Failed to delete blueprint"
)
