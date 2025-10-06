package handler

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type ValidationError struct {
	Key       string `json:"key"`
	Condition string `json:"condition"`
	Param  map[string]string  `json:"param,omitempty"`
}

func (h *Handler) ValidateRequestBody(c echo.Context, payload interface{}) []*ValidationError {
	validate := validator.New(validator.WithRequiredStructEnabled())
	var errors []*ValidationError
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
			param := validationErr.Param()
			paramMap := make(map[string]string)
			switch condition {
			case "min":
				paramMap["requiredLength"]  = param
			case "max":
				paramMap["maximumLength"] = param
			}

			currentValidationError := ValidationError{
				Key:       key,
				Param: paramMap,
				Condition: condition,
			}
			errors = append(errors, &currentValidationError)
		}
	}

	return errors
}
