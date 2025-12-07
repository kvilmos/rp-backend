package handler

import (
	"net/http"
	"room-planner/app"
	"room-planner/dto"
	"room-planner/response"

	"github.com/labstack/echo/v4"
)

func (h *Handler) HandleGetFurnitureCategory(c echo.Context) error {
	ctx := c.Request().Context()
	categories, err := h.FurnitureCatService.GetCategories(ctx)
	if err != nil {
		return NewApiError(http.StatusInternalServerError, ERROR_RETRIEVING_FURNITURE_CATEGORIES, err)
	}

	categoriesDto := dto.FromFurnitureCategoryModels(categories)

	return response.SendSuccessResponse(c, app.FURNITURE_CATEGORIES_RETRIEVED, categoriesDto)
}
