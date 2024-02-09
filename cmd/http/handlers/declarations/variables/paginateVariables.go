package variables

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations/variables"
	"creatif/pkg/app/auth"
	paginateVariables2 "creatif/pkg/app/services/variables/paginateVariables"
	"creatif/pkg/lib/logger"
	"creatif/pkg/lib/sdk"
	"github.com/labstack/echo/v4"
	"net/http"
)

func PaginateVariablesHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model variables.PaginateVariables
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = variables.SanitizePaginateVariables(model)

		l := logger.NewLogBuilder()
		handler := paginateVariables2.New(paginateVariables2.NewModel(
			model.ProjectID,
			model.Name,
			model.SanitizedLocales,
			model.OrderBy,
			model.Search,
			model.OrderDirection,
			model.Limit,
			model.Page,
			model.SanitizedGroups,
			model.Behaviour,
			sdk.ParseFilters(model.Filters),
		), auth.NewNoopAuthentication(), l)

		return request.SendResponse[paginateVariables2.Model](handler, c, http.StatusOK, l, nil, false)
	}
}
