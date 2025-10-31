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

func (h Handler) HandleCreateBlueprint(c echo.Context) error {
	claims, ok := c.Get("user_claims").(token.UserClaims)
	if !ok {
		return NewApiError(http.StatusUnauthorized, UNAUTHORIZED_REQUEST, app.ErrClaimsParsingFailed)
	}

	blueprintReq := new(request.NewBlueprintRequest)
	blueprintReq.UserId = claims.Id
	ctx := context.Background()
	bp, err := h.BpService.CreateBlueprint(ctx, *blueprintReq)
	if err != nil {
		return err
	}
	bpDto := dto.FromBlueprintModel(*bp)

	return response.SendSuccessResponse(c, "Blueprint created", bpDto)
}

func (h Handler) HandleSaveBlueprint(c echo.Context) error {
	claims, ok := c.Get("user_claims").(token.UserClaims)
	if !ok {
		return NewApiError(http.StatusUnauthorized, UNAUTHORIZED_REQUEST, app.ErrClaimsParsingFailed)
	}

	blueprintReq := new(request.NewBlueprintRequest)
	err := (&echo.DefaultBinder{}).BindBody(c, blueprintReq)
	if err != nil {
		return NewApiError(http.StatusBadRequest, INVALID_PAYLOAD, err)
	}

	pageStr := c.Param("id")
	blueprintId := int(blueprintReq.Id)
	if pageStr != "" {
		blueprintId, err = strconv.Atoi(pageStr)
		if err != nil {
			return err
		}
	}

	ctx := context.Background()
	blueprint, err := h.BpService.GetBlueprintById(ctx, int64(blueprintId))
	if err != nil {
		return err
	}
	if blueprint.Id != blueprintReq.Id {
		return NewApiError(http.StatusBadRequest, "", app.ErrIdMismatch)
	}

	if blueprint.UserId == 0 {
		blueprint, err = h.BpService.CreateBlueprint(ctx, *blueprintReq)
		if err != nil {
			return err
		}
	}

	if blueprint.UserId != claims.Id {
		return NewApiError(http.StatusUnauthorized, UPDATE_BAD_REQUEST, app.ErrCreatorIdMismatch)
	}

	blueprintReq.Id = blueprint.Id
	blueprintReq.UserId = claims.Id
	blueprint, err = h.BpService.SaveBlueprint(ctx, *blueprintReq)
	if err != nil {
		return err
	}

	bpDto := dto.FromBlueprintModel(*blueprint)

	return response.SendSuccessResponse(c, "Blueprint saved", bpDto)
}

func (h Handler) HandleGetCompleteBlueprintById(c echo.Context) error {
	claims, ok := c.Get("user_claims").(token.UserClaims)
	if !ok {
		return NewApiError(http.StatusUnauthorized, UNAUTHORIZED_REQUEST, app.ErrClaimsParsingFailed)
	}

	pageStr := c.Param("id")
	id, err := strconv.Atoi(pageStr)
	if err != nil {
		return err
	}

	ctx := context.Background()
	blueprint, err := h.BpService.GetCompleteBlueprintById(ctx, claims.Id, int64(id))
	if err != nil {
		return err
	}

	processedFurniture := make(map[int64]bool)
	var furnitureListDto []dto.FurnitureDto
	for _, item := range blueprint.Items {
		if processedFurniture[item.FurnitureId] {
			continue
		}
		furnitureDto := dto.FromFurnitureModel(item.Furniture)

		thumbnailUrl, err := h.FurnitureService.GetFurnitureFileUrl(ctx, item.Furniture.FileName, constant.THUMBNAIL_BUCKET)
		if err != nil {
			return err
		}

		objectUrl, err := h.FurnitureService.GetFurnitureFileUrl(ctx, item.Furniture.FileName, constant.FURNITURE_BUCKET)
		if err != nil {
			return err
		}
		furnitureDto.ThumbnailUrl = thumbnailUrl.String()
		furnitureDto.ObjectUrl = objectUrl.String()

		furnitureListDto = append(furnitureListDto, *furnitureDto)
		processedFurniture[item.FurnitureId] = true
	}

	blueprintDto := dto.FromCompleteBlueprintModel(*blueprint)
	blueprintDto.Furniture = furnitureListDto

	return response.SendSuccessResponse(c, "Complete Blueprint ", blueprintDto)
}

func (h Handler) HandlerGetProfileBlueprints(c echo.Context) error {
	claims, ok := c.Get("user_claims").(token.UserClaims)
	if !ok {
		return NewApiError(http.StatusUnauthorized, UNAUTHORIZED_REQUEST, app.ErrClaimsParsingFailed)
	}

	pageStr := c.QueryParam(string(constant.PAGE))
	if pageStr == "" {
		pageStr = "1"
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}
	order := c.QueryParam(string(constant.ORDER))
	if order == "" {
		order = string(constant.RECENTLY_MODIFIED)
	}
	filter := request.BlueprintFilter{
		CreatorId: claims.Id,
		Page:      page,
		Order:     order,
	}

	ctx := c.Request().Context()
	blueprints, totalRows, err := h.BpService.PageForUser(ctx, filter)
	if err != nil {
		return NewApiError(http.StatusInternalServerError, ERROR_RETRIEVING_USER_BLUEPRINTS, err)
	}
	totalPages := math.Ceil(float64(totalRows) / constant.PAGE_LIMIT)
	bpDto := dto.FromBlueprintModels(blueprints)

	pageDto := dto.BlueprintPageDto{
		CurrPage:   page,
		TotalPages: int(totalPages),
		List:       bpDto,
	}

	return response.SendSuccessResponse(c, app.USER_BLUEPRINTS_RETRIEVED, pageDto)
}
