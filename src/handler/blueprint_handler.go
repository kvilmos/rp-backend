package handler

import (
	"context"
	"net/http"
	"room-planner/app"
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
	bpId := int(blueprintReq.Id)
	if pageStr != "" {
		bpId, err = strconv.Atoi(pageStr)
		if err != nil {
			return err
		}
	}

	ctx := context.Background()
	bp, err := h.BpService.GetBlueprintById(ctx, int64(bpId))
	if err != nil {
		return err
	}

	if bp.UserId == 0 {
		bp, err = h.BpService.CreateBlueprint(ctx, *blueprintReq)
		if err != nil {
			return err
		}
	}

	if bp.UserId != claims.Id {
		return NewApiError(http.StatusUnauthorized, UNAUTHORIZED_REQUEST, app.ErrCreatorIdMismatch)
	}

	blueprintReq.Id = bp.Id
	blueprintReq.UserId = claims.Id
	bp, err = h.BpService.SaveBlueprint(ctx, *blueprintReq)
	if err != nil {
		return err
	}

	bpDto := dto.FromBlueprintModel(*bp)

	return response.SendSuccessResponse(c, "Blueprint created", bpDto)
}

func (h Handler) HandleListBlueprints(c echo.Context) error {

	return response.SendSuccessResponse(c, "Blueprint created", nil)
}

func (h Handler) HandleGetBlueprintById(c echo.Context) error {

	return response.SendSuccessResponse(c, "Blueprint created", nil)
}
