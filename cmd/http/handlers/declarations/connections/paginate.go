package connections

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations/connections"
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/connections/pagination"
	"creatif/pkg/lib/sdk"
	"github.com/labstack/echo/v4"
	"net/http"
)

func PaginateConnectionsHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model connections.PaginateConnections
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = connections.SanitizePaginateConnections(model)

		authentication := auth.NewApiAuthentication(request.GetApiAuthenticationCookie(c))
		handler := pagination.New(pagination.NewModel(
			model.ProjectID,
			model.StructureType,
			model.ParentVariableID,
			model.SanitizedLocales,
			model.StructureID,
			model.OrderBy,
			model.Search,
			model.OrderDirection,
			model.Limit,
			model.Page,
			model.SanitizedGroups,
			sdk.ParseFilters(model.Filters),
			model.Behaviour,
			model.SanitizedFields,
		), authentication)

		return request.SendResponse[pagination.Model](handler, c, http.StatusOK, func(c echo.Context, model interface{}) error {
			if authentication.ShouldRefresh() {
				session, err := authentication.Refresh()
				if err != nil {
					return err
				}

				c.SetCookie(request.EncryptAuthenticationCookie(session))
			}

			return nil
		}, false)
	}
}
