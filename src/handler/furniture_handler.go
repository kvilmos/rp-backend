package handler

import (
	"context"
	"math"
	"net/http"
	"room-planner/app"
	"room-planner/common/constant"
	"room-planner/dto"
	"room-planner/request"
	"room-planner/response"
	"room-planner/token"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (h *Handler) HandleNewFurniture(c echo.Context) error {
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

func (h *Handler) HandleGetFurnitureList(c echo.Context) error {
	filter := request.FurnitureFilter{}

	pageStr := c.QueryParam(string(constant.PAGE))
	if pageStr == "" {
		pageStr = "1"
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}
	filter.Page = page

	order := c.QueryParam(string(constant.ORDER))
	if order == "" {
		order = string(constant.RECENTLY_MODIFIED)
	}
	filter.Order = order

	category := c.QueryParam(string(constant.CATEGORY_ID))
	categoryId, err := strconv.ParseInt(category, 10, 64)
	if err != nil {
		filter.CategoryId = nil
	} else {
		filter.CategoryId = &categoryId
	}

	ctx := c.Request().Context()
	furnitureList, totalRows, err := h.FurnitureService.PageForUser(ctx, filter)
	if err != nil {
		return NewApiError(http.StatusInternalServerError, ERROR_RETRIEVING_USER_FURNITURE_LIST, err)
	}
	totalPages := math.Ceil(float64(totalRows) / constant.PAGE_LIMIT)

	var furnitureListDto []dto.FurnitureDto
	for _, furniture := range furnitureList {
		furnitureDto := dto.FromFurnitureModel(furniture)

		thumbnailUrl, err := h.FurnitureService.GetFurnitureFileUrl(ctx, furniture.FileName, constant.THUMBNAIL_BUCKET)
		if err != nil {
			return err
		}

		objectUrl, err := h.FurnitureService.GetFurnitureFileUrl(ctx, furniture.FileName, constant.FURNITURE_BUCKET)
		if err != nil {
			return err
		}

		furnitureDto.ObjectUrl = objectUrl.String()
		furnitureDto.ThumbnailUrl = thumbnailUrl.String()

		furnitureListDto = append(furnitureListDto, *furnitureDto)
	}

	responseDto := dto.FurniturePageDto{
		CurrPage:   page,
		TotalPages: int(totalPages),
		List:       furnitureListDto,
	}

	return response.SendSuccessResponse(c, app.USER_FURNITURE_LIST_RETRIEVED, responseDto)
}

func (h *Handler) HandlerGetProfileFurniture(c echo.Context) error {
	filter := request.FurnitureFilter{}

	claims, ok := c.Get("user_claims").(token.UserClaims)
	if !ok {
		return NewApiError(http.StatusUnauthorized, UNAUTHORIZED_REQUEST, app.ErrClaimsParsingFailed)
	}
	filter.CreatorId = claims.Id

	pageStr := c.QueryParam(string(constant.PAGE))
	if pageStr == "" {
		pageStr = "1"
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}
	filter.Page = page

	order := c.QueryParam(string(constant.ORDER))
	if order == "" {
		order = string(constant.RECENTLY_MODIFIED)
	}
	filter.Order = order

	category := c.QueryParam(string(constant.CATEGORY_ID))
	categoryId, err := strconv.ParseInt(category, 10, 64)
	if err != nil {
		filter.CategoryId = nil
	} else {
		filter.CategoryId = &categoryId
	}

	ctx := c.Request().Context()
	furnitureList, totalRows, err := h.FurnitureService.PageForUser(ctx, filter)
	if err != nil {
		return NewApiError(http.StatusInternalServerError, ERROR_RETRIEVING_USER_FURNITURE_LIST, err)
	}
	totalPages := math.Ceil(float64(totalRows) / constant.PAGE_LIMIT)

	var furnitureListDto []dto.FurnitureDto
	for _, furniture := range furnitureList {
		furnitureDto := dto.FromFurnitureModel(furniture)

		thumbnailUrl, err := h.FurnitureService.GetFurnitureFileUrl(ctx, furniture.FileName, constant.THUMBNAIL_BUCKET)
		if err != nil {
			return err
		}

		objectUrl, err := h.FurnitureService.GetFurnitureFileUrl(ctx, furniture.FileName, constant.FURNITURE_BUCKET)
		if err != nil {
			return err
		}

		furnitureDto.ObjectUrl = objectUrl.String()
		furnitureDto.ThumbnailUrl = thumbnailUrl.String()

		furnitureListDto = append(furnitureListDto, *furnitureDto)
	}

	responseDto := dto.FurniturePageDto{
		CurrPage:   page,
		TotalPages: int(totalPages),
		List:       furnitureListDto,
	}

	return response.SendSuccessResponse(c, app.USER_FURNITURE_LIST_RETRIEVED, responseDto)
}

func (h *Handler) HandleGetFurnitureById(c echo.Context) error {
	pageStr := c.Param("id")
	id, err := strconv.Atoi(pageStr)
	if err != nil {
		return err
	}
	ctx := context.Background()

	furniture, err := h.FurnitureService.GetFurnitureById(ctx, int64(id))
	if err != nil {
		return err
	}

	furnDto := dto.FromFurnitureModel(furniture)

	thumbnailUrl, err := h.FurnitureService.GetFurnitureFileUrl(ctx, furniture.FileName, constant.THUMBNAIL_BUCKET)
	if err != nil {
		return err
	}

	objectUrl, err := h.FurnitureService.GetFurnitureFileUrl(ctx, furniture.FileName, constant.FURNITURE_BUCKET)
	if err != nil {
		return err
	}

	furnDto.ThumbnailUrl = thumbnailUrl.String()
	furnDto.ObjectUrl = objectUrl.String()

	return response.SendSuccessResponse(c, "furniture", furnDto)
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
