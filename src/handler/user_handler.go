package handler

import (
	"room-planner/model"
	"room-planner/request"
	"room-planner/response"

	"github.com/labstack/echo/v4"
)

func (h Handler) Register(c echo.Context) error {
	regReq := new(request.RegisterRequest)
	err := (&echo.DefaultBinder{}).BindBody(c, regReq)
	if err != nil {
		return response.SendInternalServerErrorResponse(c, err.Error())
	}

	validationErrors := h.ValidateRequestBody(c, *regReq)
	if validationErrors != nil {
		return response.SendFailedValidationResponse(c, validationErrors)
	}

	user, err := h.UserService.RegisterUser(*regReq)
	if err != nil {
		if err.Error() == "email has already been taken" {
			return response.SendBadRequestResponse(c, err.Error())
		}
		return response.SendInternalServerErrorResponse(c, err.Error())
	}

	return response.SendSuccessResponse(c, "Sign up successful", user)
}

func (h Handler) Login(c echo.Context) error {
	loginReq := new(request.LoginRequest)
	err := (&echo.DefaultBinder{}).BindBody(c, loginReq)
	if err != nil {
		return response.SendBadRequestResponse(c, err.Error())
	}

	validationErrors := h.ValidateRequestBody(c, *loginReq)
	if validationErrors != nil {
		return response.SendFailedValidationResponse(c, validationErrors)
	}

	accessToken, refreshToken, userRetrieved, err := h.UserService.LoginUser(*loginReq)
	if err != nil {
		return response.SendBadRequestResponse(c, err.Error())
	}

	return response.SendSuccessResponse(c, "User logged in", map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user":          userRetrieved,
	})
}

func (h Handler) TestAuth(c echo.Context) error {
	user, ok := c.Get("user").(model.User)
	if !ok {
		return response.SendInternalServerErrorResponse(c, "User authentication failed")
	}

	return response.SendSuccessResponse(c, "Authenticated user retrieved", user)
}
