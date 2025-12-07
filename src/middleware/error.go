package middleware

import (
	"log/slog"
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

		slog.Error(err.Error())

		apiError, ok := err.(handler.ApiError)
		if ok {
			return response.SendErrorResponse(c, apiError.Status, apiError.Msg)
		}

		return response.SendErrorResponse(c, http.StatusInternalServerError, "internalServerError")
	}
}
