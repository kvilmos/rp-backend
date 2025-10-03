package handler

import (
	"fmt"
	"reflect"
	"room-planner/response"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func (h *Handler) ValidateRequestBody(c echo.Context, payload interface{}) []*response.ValidationError {
	validate := validator.New(validator.WithRequiredStructEnabled())
	var errors []*response.ValidationError
	err := validate.Struct(payload)
	validationErrors, ok := err.(validator.ValidationErrors)
	if ok {
		reflected := reflect.ValueOf(payload)

		for _, validationErr := range validationErrors {
			field, _ := reflected.Type().FieldByName(validationErr.StructField())

			key := field.Tag.Get("json")
			if key == "" {
				key = strings.ToLower(validationErr.StructField())
			}

			condition := validationErr.Tag()
			keyToTitleCase := strings.Replace(key, "_", " ", -1)
			errMessage := keyToTitleCase + " field is " + condition
			param := validationErr.Param()

			switch condition {
			case "required":
				errMessage = keyToTitleCase + " required"
			case "email":
				errMessage = keyToTitleCase + " must be a valid email"
			case "min":
				errMessage = fmt.Sprintf("%s must be az least %s character", keyToTitleCase, param)
			case "max":
				errMessage = fmt.Sprintf("%s cannot be az more than %s character", keyToTitleCase, param)
			}

			currentValidationError := response.ValidationError{
				Error:     errMessage,
				Key:       key,
				Condition: condition,
			}
			errors = append(errors, &currentValidationError)
		}
	}

	return errors
}
