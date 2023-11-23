package variables

import (
	declarations2 "creatif/cmd/http/handlers/declarations"
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
		if model.Locale == "" {
			model.Locale = declarations2.DefaultLocale
		}

		l := logger.NewLogBuilder()
		handler := paginateVariables2.New(paginateVariables2.NewModel(
			model.ProjectID,
			model.Locale,
			model.OrderBy,
			model.Search,
			model.OrderDirection,
			model.Limit,
			model.Page,
			model.SanitizedGroups,
			sdk.ParseFilters(model.Filters),
		), auth.NewNoopAuthentication(), l)

		return request.SendResponse[paginateVariables2.Model](handler, c, http.StatusOK, l, nil, false)
	}
}
