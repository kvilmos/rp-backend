package middleware

import (
	"room-planner/response"
	"room-planner/service"
	"room-planner/util"
	"strings"

	"github.com/labstack/echo/v4"
)

type AuthMiddleware struct {
	userService *service.UserService
}

func NewAuthMiddleware(us *service.UserService) *AuthMiddleware {
	return &AuthMiddleware{
		userService: us,
	}
}

func (appMiddleware *AuthMiddleware) Authenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Add("Vary", "Authorization")
		authHeader := c.Request().Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer") {
			return response.SendUnauthorizedResponse(c, "Please provide a Bearer token")
		}

		authHeaderSplits := strings.Split(authHeader, " ")
		accessToken := authHeaderSplits[1]

		claims, err := util.ParseJWTSignedAccessToken(accessToken)
		if err != nil {
			return response.SendUnauthorizedResponse(c, err.Error())
		}

		if util.IsClaimExpired(claims) {
			return response.SendUnauthorizedResponse(c, "Token is expired")
		}

		user, err := appMiddleware.userService.GetUserById(claims.Id)
		if err != nil {
			return response.SendInternalServerErrorResponse(c, err.Error())
		}

		c.Set("user", user)

		return next(c)
	}
}
