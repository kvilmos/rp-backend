package handler

import (
	"errors"
	"net/http"
	"room-planner/app"
	"room-planner/common/constant"
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
		return NewApiError(http.StatusBadRequest, INVALID_PAYLOAD, err)
	}

	validationErrors := h.ValidateRequestBody(c, *regReq)
	if validationErrors != nil {
		return NewApiError(http.StatusUnprocessableEntity, validationErrors, app.ErrValidationFailed)
	}

	user, err := h.UserService.RegisterUser(*regReq)
	if err != nil {
		if errors.Is(err, app.ErrAlreadyExist) {
			return NewApiError(http.StatusConflict, EMAIL_ALREADY_EXIST, err)
		}
		return err
	}
	userDto := dto.FromUserModel(user)

	return response.SendSuccessResponse(c, app.REGISTRATION_SUCCESS, userDto)
}

func (h Handler) HandlerLoginUser(c echo.Context) error {
	loginReq := new(request.LoginRequest)
	err := (&echo.DefaultBinder{}).BindBody(c, loginReq)
	if err != nil {
		return NewApiError(http.StatusBadRequest, INVALID_PAYLOAD, err)
	}

	validationErrors := h.ValidateRequestBody(c, *loginReq)
	if validationErrors != nil {
		return NewApiError(http.StatusUnprocessableEntity, validationErrors, app.ErrValidationFailed)
	}

	accessToken, refreshToken, refreshClaims, user, err := h.UserService.LoginUser(*loginReq)
	if err != nil {
		if errors.Is(err, app.ErrInvalidCredentials) {
			return NewApiError(http.StatusConflict, INVALID_LOGIN_CREDENTIALS, err)
		}
		return err
	}

	expiresAt := refreshClaims.RegisteredClaims.ExpiresAt.Time
	maxAge := int(time.Until(expiresAt).Seconds())
	cookie := &http.Cookie{
		Name:     constant.REFRESH_TOKEN_COOKIE,
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

	loginDto := dto.LoginDto{
		AccessToken: *accessToken,
		User:        *dto.FromUserModel(user),
	}

	return response.SendSuccessResponse(c, app.LOGIN_SUCCESS, loginDto)
}

func (h Handler) HandlerRenewToken(c echo.Context) error {
	cookie, err := c.Cookie(constant.REFRESH_TOKEN_COOKIE)
	if err != nil {
		return NewApiError(http.StatusUnauthorized, UNAUTHORIZED_REQUEST, err)
	}
	refreshToken := cookie.Value

	newAccessToken, newRefreshToken, newRefreshClaims, err := h.UserService.RenewUserAccessToken(refreshToken)
	if err != nil {
		return NewApiError(http.StatusUnauthorized, UNAUTHORIZED_REQUEST, err)
	}

	expiresAt := newRefreshClaims.RegisteredClaims.ExpiresAt.Time
	maxAge := int(time.Until(expiresAt).Seconds())
	newCookie := &http.Cookie{
		Name:     constant.REFRESH_TOKEN_COOKIE,
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

	return response.SendSuccessResponse(c, app.TOKEN_RENEWED, dto.RenewAccessTokenDto{
		AccessToken: *newAccessToken,
	})
}

func (h Handler) HandleLogoutUser(c echo.Context) error {
	cookie, err := c.Cookie(constant.REFRESH_TOKEN_COOKIE)
	if err != nil {
		return err
	}

	if cookie != nil {
		refreshToken := cookie.Value
		err = h.UserService.Logout(refreshToken)
		if err != nil {
			return err
		}
	}

	newCookie := &http.Cookie{
		Name:     constant.REFRESH_TOKEN_COOKIE,
		Value:    "",
		Expires:  time.Unix(0, 0),
		MaxAge:   0,
		Path:     "/",
		Domain:   "localhost",
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	c.SetCookie(newCookie)

	return response.SendSuccessResponse(c, app.LOGOUT_SUCCESS, "")
}

func (h Handler) HandleVerifyUser(c echo.Context) error {
	claims, ok := c.Get(constant.USER_CLAIMS).(token.UserClaims)
	if !ok {
		return NewApiError(http.StatusUnauthorized, UNAUTHORIZED_REQUEST, app.ErrClaimsParsingFailed)
	}

	user, err := h.UserService.GetUserById(claims.Id)
	if err != nil {
		return NewApiError(http.StatusUnauthorized, UNAUTHORIZED_REQUEST, err)
	}
	userDto := dto.FromUserModel(user)

	return response.SendSuccessResponse(c, app.USER_RETRIEVED, userDto)
}
