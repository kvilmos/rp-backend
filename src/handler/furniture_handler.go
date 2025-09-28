package handler

import (
	"context"
	"fmt"
	"net/http"
	"room-planner/request"
	"room-planner/response"

	"github.com/labstack/echo/v4"
)

func (h *Handler) NewFurniture(c echo.Context) error {
	furnitureReq := new(request.NewFurnitureRequest)
	err := (&echo.DefaultBinder{}).BindBody(c, furnitureReq)
	if err != nil {
		return response.SendInternalServerErrorResponse(c, err.Error())
	}

	// MISSING VALIDATION

	ctx := context.Background()
	urls, err := h.FurnitureService.PrepareUpload(ctx, furnitureReq.Name)
	if err != nil {
		if err.Error() == "email has already been taken" {
			return response.SendBadRequestResponse(c, err.Error())
		}
		return response.SendInternalServerErrorResponse(c, err.Error())
	}

	return c.JSON(http.StatusAccepted, urls)
	//return http.resonse.SendSuccessResponse(c, "Sign up successful", urls)
}

func (h *Handler) HandleUploadNotification(c echo.Context) error {
	notification := new(request.MinioNotification)
	err := (&echo.DefaultBinder{}).BindBody(c, notification)
	if err != nil {
		return response.SendInternalServerErrorResponse(c, err.Error())
	}
	fmt.Println("HOOK NOTIFICATION: ", *notification)

	ctx := context.Background()
	err = h.FurnitureService.QueueUploadNotification(ctx, *notification)
	if err != nil {
		return response.SendInternalServerErrorResponse(c, err.Error())
	}

	return c.NoContent(http.StatusAccepted)
}
