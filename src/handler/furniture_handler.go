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

func (h *Handler) ListFurniture(c echo.Context) error {
	ctx := context.Background()
	furnitureList, err := h.FurnitureService.ListFurniture(ctx)
	if err != nil {
		return err
	}

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

	return response.SendSuccessResponse(c, "furniture list", furnitureListDto)
}

func (h *Handler) PageFurniture(c echo.Context) error {
	pageStr := c.Param("page")

	page := 1
	if pageStr != "" {
		newPage, err := strconv.Atoi(pageStr)
		if err != nil {
			return err
		}
		page = newPage
	}

	orderBy := c.QueryParam("orderBy")
	creatorIdStr := c.QueryParam("creator")

	creatorId := int64(0)
	if creatorIdStr != "" {
		IdInt, err := strconv.Atoi(creatorIdStr)
		if err != nil {
			return err
		}
		creatorId = int64(IdInt)
	}

	filter := request.FurnitureFilter{
		OrderBy:   orderBy,
		CreatorId: creatorId,
	}

	ctx := context.Background()
	furnitureList, err := h.FurnitureService.PageFurniture(ctx, page, filter)
	if err != nil {
		return err
	}

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

	var totalRows int
	totalRows, err = h.FurnitureService.GetFurnitureCount(ctx)
	if err != nil {
		return err
	}
	totalPages := math.Ceil(float64(totalRows) / constant.PAGE_LIMIT)

	responseDto := dto.FurniturePageDto{
		NextPage:   page + 1,
		PrevPage:   page - 1,
		CurrPage:   page,
		TotalPages: int(totalPages),
		List:       furnitureListDto,
	}

	return response.SendSuccessResponse(c, "furniture page", responseDto)
}

func (h *Handler) HandlerPageUserFurniture(c echo.Context) error {
	claims, ok := c.Get("user_claims").(token.UserClaims)
	if !ok {
		return NewApiError(http.StatusUnauthorized, UNAUTHORIZED_REQUEST, app.ErrClaimsParsingFailed)
	}

	userIdStr := c.Param("userId")
	userId := int64(0)
	if userIdStr != "" {
		userIdInt, err := strconv.Atoi(userIdStr)
		if err != nil {
			return err
		}
		userId = int64(userIdInt)
	}

	if claims.Id != userId {
		return NewApiError(http.StatusUnauthorized, UNAUTHORIZED_REQUEST, app.ErrClaimsParsingFailed)
	}

	pageStr := c.Param("page")
	page := 1
	if pageStr != "" {
		newPage, err := strconv.Atoi(pageStr)
		if err != nil {
			return err
		}
		page = newPage
	}

	orderBy := c.QueryParam("orderBy")
	filter := request.FurnitureFilter{
		OrderBy:   orderBy,
		CreatorId: claims.Id,
	}

	ctx := context.Background()
	furnitureList, err := h.FurnitureService.PageFurniture(ctx, page, filter)
	if err != nil {
		return err
	}

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

	var totalRows int
	totalRows, err = h.FurnitureService.GetFurnitureCount(ctx)
	if err != nil {
		return err
	}
	totalPages := math.Ceil(float64(totalRows) / constant.PAGE_LIMIT)

	responseDto := dto.FurniturePageDto{
		NextPage:   page + 1,
		PrevPage:   page - 1,
		CurrPage:   page,
		TotalPages: int(totalPages),
		List:       furnitureListDto,
	}

	return response.SendSuccessResponse(c, "furniture page", responseDto)
}

func (h *Handler) GetFurnitureById(c echo.Context) error {
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
