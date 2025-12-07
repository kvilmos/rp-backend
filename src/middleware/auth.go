package middleware

import (
	"net/http"
	"room-planner/app"
	"room-planner/common/constant"
	"room-planner/handler"
	"room-planner/service"
	"room-planner/token"
	"strings"

	"github.com/labstack/echo/v4"
)

type AuthMiddleware struct {
	userService *service.UserService
	JWTMaker    *token.JWTMaker
}

func NewAuthMiddleware(us *service.UserService, jwtMaker *token.JWTMaker) *AuthMiddleware {
	return &AuthMiddleware{
		userService: us,
		JWTMaker:    jwtMaker,
	}
}

func (auth *AuthMiddleware) Authenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Add("Vary", "Authorization")
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return handler.NewApiError(http.StatusUnauthorized, handler.UNAUTHORIZED_REQUEST, app.ErrAuthHeaderMissing)
		}

		fields := strings.Fields(authHeader)
		if len(fields) != 2 || fields[0] != "Bearer" {
			return handler.NewApiError(http.StatusUnauthorized, handler.UNAUTHORIZED_REQUEST, app.ErrInvalidAuthHeader)
		}

		token := fields[1]
		claims, err := auth.JWTMaker.VerifyToken(token)
		if err != nil {
			return handler.NewApiError(http.StatusUnauthorized, handler.UNAUTHORIZED_REQUEST, err)
		}
		c.Set(constant.USER_CLAIMS, *claims)

		return next(c)
	}
}
