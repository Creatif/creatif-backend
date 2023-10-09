package declarations

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations"
	paginateVariables2 "creatif/pkg/app/services/variables/paginateVariables"
	"creatif/pkg/lib/sdk"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

func PaginateVariablesHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model declarations.PaginateVariables
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = declarations.SanitizePaginateVariables(model)
		model.OrderBy = strings.ToUpper(model.OrderBy)

		handler := paginateVariables2.New(paginateVariables2.NewModel(
			model.ProjectID,
			model.OrderBy,
			model.OrderDirection,
			model.Limit,
			model.Page,
			model.SanitizedGroups,
			sdk.ParseFilters(model.Filters),
		))

		return request.SendResponse[paginateVariables2.Model](handler, c, http.StatusOK)
	}
}
