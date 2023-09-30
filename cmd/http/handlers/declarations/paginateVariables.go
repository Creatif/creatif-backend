package declarations

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations"
	"creatif/pkg/app/declarations/paginateVariables"
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

		handler := paginateVariables.New(paginateVariables.NewModel(model.ProjectID, model.PaginationID, model.Field, model.OrderBy, model.Direction, model.Limit, model.SanitizedGroups))

		return request.SendResponse[paginateVariables.Model](handler, c, http.StatusOK)
	}
}
