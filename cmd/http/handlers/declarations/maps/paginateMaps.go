package maps

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations/maps"
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/maps/paginateMaps"
	"github.com/labstack/echo/v4"
	"net/http"
)

func PaginateMapsHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model maps.PaginateMaps
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = maps.SanitizePaginateMaps(model)

		authentication := auth.NewApiAuthentication(request.GetApiAuthenticationCookie(c))
		handler := paginateMaps.New(paginateMaps.NewModel(
			model.ProjectID,
			model.OrderBy,
			model.Search,
			model.OrderDirection,
			model.Limit,
			model.Page,
		), authentication)

		return request.SendResponse[paginateMaps.Model](handler, c, http.StatusOK, func(c echo.Context, model interface{}) error {
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
