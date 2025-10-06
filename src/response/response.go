package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type JSONSuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type JSONStatusResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type JSONErrorResponse struct {
	Success bool `json:"success"`
	Message any  `json:"message"`
}

func SendSuccessResponse(c echo.Context, message string, data interface{}) error {
	return c.JSON(http.StatusOK, JSONSuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func SendStatusResponse(c echo.Context, statusCode int, message string) error {
	return c.JSON(statusCode, JSONStatusResponse{
		Success: true,
		Message: message,
	})
}

func SendErrorResponse(c echo.Context, statusCode int, message any) error {
	return c.JSON(statusCode, JSONErrorResponse{
		Success: false,
		Message: message,
	})
}
