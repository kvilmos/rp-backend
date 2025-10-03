package middleware

import (
	"fmt"
	"room-planner/response"
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
			return response.SendUnauthorizedResponse(c, "authorization header is missing")
		}

		fields := strings.Fields(authHeader)
		if len(fields) != 2 || fields[0] != "Bearer" {
			return response.SendUnauthorizedResponse(c, "invalid authorization header")
		}

		token := fields[1]
		claims, err := auth.JWTMaker.VerifyToken(token)
		if err != nil {
			fmt.Println(err)
			return response.SendUnauthorizedResponse(c, err.Error())
		}
		c.Set("user_claims", *claims)

		return next(c)
	}
}
