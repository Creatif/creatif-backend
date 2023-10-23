package lists

import (
	declarations2 "creatif/cmd/http/handlers/declarations"
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations/lists"
	"creatif/pkg/app/services/lists/paginateListItems"
	"creatif/pkg/lib/sdk"
	"github.com/labstack/echo/v4"
	"net/http"
)

func PaginateListItemsHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model lists.PaginateListItems
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = lists.SanitizePaginateListItems(model)
		if model.Locale == "" {
			model.Locale = declarations2.DefaultLocale
		}

		handler := paginateListItems.New(paginateListItems.NewModel(
			model.ProjectID,
			model.Locale,
			model.ListName,
			model.OrderBy,
			model.OrderDirection,
			model.Limit,
			model.Page,
			model.SanitizedGroups,
			sdk.ParseFilters(model.Filters),
		))

		return request.SendResponse[paginateListItems.Model](handler, c, http.StatusOK)
	}
}
