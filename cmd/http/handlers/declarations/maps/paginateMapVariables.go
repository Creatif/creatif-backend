package maps

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations/maps"
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/maps/paginateMapVariables"
	"creatif/pkg/lib/logger"
	"creatif/pkg/lib/sdk"
	"github.com/labstack/echo/v4"
	"net/http"
)

func PaginateMapVariables() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model maps.PaginateMapVariables
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = maps.SanitizePaginateListItems(model)

		l := logger.NewLogBuilder()
		authentication := auth.NewApiAuthentication(request.GetApiAuthenticationCookie(c), l)
		handler := paginateMapVariables.New(paginateMapVariables.NewModel(
			model.ProjectID,
			model.SanitizedLocales,
			model.ListName,
			model.OrderBy,
			model.Search,
			model.OrderDirection,
			model.Limit,
			model.Page,
			model.SanitizedGroups,
			sdk.ParseFilters(model.Filters),
			model.Behaviour,
			model.SanitizedFields,
		), authentication, l)

		return request.SendResponse[paginateMapVariables.Model](handler, c, http.StatusOK, l, func(c echo.Context, model interface{}) error {
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
