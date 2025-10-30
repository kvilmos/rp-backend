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

func (h Handler) HandleListBlueprints(c echo.Context) error {
	ctx := context.Background()
	blueprints, err := h.BpService.ListBlueprints(ctx)
	if err != nil {
		return err
	}
	blueprintsDto := dto.FromBlueprintsModel(blueprints)

	return response.SendSuccessResponse(c, "All Blueprints listed", blueprintsDto)
}

func (h Handler) HandlePageBlueprints(c echo.Context) error {
	pageStr := c.Param("page")

	page := 1
	if pageStr != "" {
		newPage, err := strconv.Atoi(pageStr)
		if err != nil {
			return err
		}
		page = newPage
	}

	ctx := context.Background()
	blueprints, err := h.BpService.PageBlueprints(ctx, page)
	if err != nil {
		return err
	}

	var totalRows int
	totalRows, err = h.BpService.GetBlueprintCount(ctx)
	if err != nil {
		return err
	}
	totalPages := math.Ceil(float64(totalRows) / constant.PAGE_LIMIT)

	pageDto := dto.BlueprintPageDto{
		NextPage:   page + 1,
		PrevPage:   page - 1,
		CurrPage:   page,
		TotalPages: int(totalPages),
		List:       dto.FromBlueprintsModel(blueprints),
	}

	return response.SendSuccessResponse(c, "Blueprint paged", pageDto)
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

func (h Handler) HandleListCompleteBlueprints(c echo.Context) error {
	ctx := context.Background()
	blueprints, err := h.BpService.GetCompleteBlueprints(ctx)
	if err != nil {
		return err
	}
	blueprintsDto := dto.FromCompleteBlueprintsModel(blueprints)

	return response.SendSuccessResponse(c, "All Complete blueprints listed", blueprintsDto)
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

func (h Handler) HandlerPageUserBlueprint(c echo.Context) error {
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

	orderBy := c.QueryParam("orderBy")

	filter := request.BlueprintFilter{
		OrderBy:   orderBy,
		CreatorId: claims.Id,
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

	ctx := context.Background()
	blueprints, err := h.BpService.GetBlueprintsByFilter(ctx, page, filter)
	if err != nil {
		return err
	}

	var totalRows int
	totalRows, err = h.BpService.GetUserBlueprintCount(ctx, claims.Id)
	if err != nil {
		return err
	}
	totalPages := math.Ceil(float64(totalRows) / constant.PAGE_LIMIT)

	pageDto := dto.BlueprintPageDto{
		NextPage:   page + 1,
		PrevPage:   page - 1,
		CurrPage:   page,
		TotalPages: int(totalPages),
		List:       dto.FromBlueprintsModel(blueprints),
	}

	return response.SendSuccessResponse(c, "All Complete blueprints listed", pageDto)
}
