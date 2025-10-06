package middleware

import (
	"fmt"
	"net/http"
	"room-planner/handler"
	"room-planner/response"

	"github.com/labstack/echo/v4"
)

func ErrorMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		if err == nil {
			return nil
		}

		// TODO Logger
		fmt.Println(err)

		apiError, ok := err.(handler.ApiError)
		if ok {
			return response.SendErrorResponse(c, apiError.Status, apiError.Msg)
		}

		return response.SendErrorResponse(c, http.StatusInternalServerError, "internalServerError")
	}
}
