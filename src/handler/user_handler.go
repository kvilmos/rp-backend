package handler

import (
	"fmt"
	"net/http"
	"room-planner/dto"
	"room-planner/request"
	"room-planner/response"
	"room-planner/token"
	"time"

	"github.com/labstack/echo/v4"
)

func (h Handler) HandlerRegisterUser(c echo.Context) error {
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

	userDto := dto.FromUserModel(user)

	return response.SendSuccessResponse(c, "Sign up successful", userDto)
}

func (h Handler) HandlerLoginUser(c echo.Context) error {
	loginReq := new(request.LoginRequest)
	err := (&echo.DefaultBinder{}).BindBody(c, loginReq)
	if err != nil {
		return response.SendBadRequestResponse(c, err.Error())
	}

	validationErrors := h.ValidateRequestBody(c, *loginReq)
	if validationErrors != nil {
		return response.SendFailedValidationResponse(c, validationErrors)
	}

	accessToken, refreshToken, refreshClaims, user, err := h.UserService.LoginUser(*loginReq)
	if err != nil {
		return response.SendBadRequestResponse(c, err.Error())
	}

	expiresAt := refreshClaims.RegisteredClaims.ExpiresAt.Time
	maxAge := int(time.Until(expiresAt).Seconds())
	cookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    *refreshToken,
		Expires:  expiresAt,
		MaxAge:   maxAge,
		Path:     "/",
		Domain:   "localhost",
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	c.SetCookie(cookie)

	return response.SendSuccessResponse(c, "User logged in", dto.LoginDTO{
		AccessToken: *accessToken,
		User:        *dto.FromUserModel(user),
	})
}

func (h Handler) HandlerRenewToken(c echo.Context) error {
	cookie, err := c.Cookie("refresh_token")
	if err != nil {
		return response.SendUnauthorizedResponse(c, err.Error())
	}
	refreshToken := cookie.Value

	newAccessToken, newRefreshToken, newRefreshClaims, err := h.UserService.RenewUserAccessToken(refreshToken)
	if err != nil {
		return response.SendInternalServerErrorResponse(c, err.Error())
	}

	expiresAt := newRefreshClaims.RegisteredClaims.ExpiresAt.Time
	maxAge := int(time.Until(expiresAt).Seconds())
	newCookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    *newRefreshToken,
		Expires:  expiresAt,
		MaxAge:   maxAge,
		Path:     "/",
		Domain:   "localhost",
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	c.SetCookie(newCookie)

	return response.SendSuccessResponse(c, "token renewed", dto.RenewAccessTokenDto{
		AccessToken: *newAccessToken,
	})
}

func (h Handler) HandleLogoutUser(c echo.Context) error {
	cookie, err := c.Cookie("refresh_token")
	if err != nil {
	}

	if cookie != nil{
		refreshToken := cookie.Value
		err = h.UserService.Logout(refreshToken)
		if err != nil {

		}
	}
	


	newCookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		Path:     "/",
		Domain:   "localhost",
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	c.SetCookie(newCookie)

	return response.SendSuccessResponse(c, "logout success", "")
}

func (h Handler) HandleVerifyUser(c echo.Context) error {
	claims, ok := c.Get("user_claims").(token.UserClaims)
	if !ok {
		return response.SendInternalServerErrorResponse(c, "User authentication failed")
	}

	user, err := h.UserService.GetUserById(claims.Id)
	if err != nil{
		return response.SendInternalServerErrorResponse(c, "User authentication failed")
	}

	userDto := dto.FromUserModel(user)

	return response.SendSuccessResponse(c, "Authenticated user retrieved", userDto)
}


func (h Handler) VerifyUser(c echo.Context) error {
	claims, ok := c.Get("user_claims").(token.UserClaims)
	fmt.Println(claims)
	if !ok {
		return response.SendInternalServerErrorResponse(c, "User authentication failed")
	}

	return response.SendSuccessResponse(c, "Authenticated user retrieved", claims)
}
