package handler

import (
	"context"
	"net/http"
	"room-planner/app"
	"room-planner/request"
	"room-planner/response"
	"room-planner/token"

	"github.com/labstack/echo/v4"
)

func (h *Handler) NewFurniture(c echo.Context) error {
	claims, ok := c.Get("user_claims").(token.UserClaims)
	if !ok {
		return NewApiError(http.StatusUnauthorized, UNAUTHORIZED_REQUEST, app.ErrClaimsParsingFailed)
	}

	furnitureReq := new(request.NewFurnitureRequest)
	err := (&echo.DefaultBinder{}).BindBody(c, furnitureReq)
	if err != nil {
		return NewApiError(http.StatusBadRequest, INVALID_PAYLOAD, err)
	}

	validationErrors := h.ValidateRequestBody(c, *furnitureReq)
	if validationErrors != nil {
		return NewApiError(http.StatusUnprocessableEntity, validationErrors, app.ErrValidationFailed)
	}

	furnitureReq.UserId = claims.Id

	ctx := context.Background()
	urls, err := h.FurnitureService.PrepareUpload(ctx, *furnitureReq)
	if err != nil {
		return err
	}

	return response.SendSuccessResponse(c, "Successful url request", urls)
}

func (h *Handler) HandleUploadNotification(c echo.Context) error {
	notification := new(request.MinioNotification)
	err := (&echo.DefaultBinder{}).BindBody(c, notification)
	if err != nil {
		return err
	}

	ctx := context.Background()
	err = h.FurnitureService.QueueUploadNotification(ctx, *notification)
	if err != nil {
		return err
	}

	return response.SendStatusResponse(c, http.StatusAccepted, "Notification accepted")
}
