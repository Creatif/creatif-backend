package maps

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations/maps"
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/maps/getMapGroups"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetMapGroupsHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model maps.GetMapGroups
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = maps.SanitizeGetMapGroups(model)

		a := auth.NewApiAuthentication(request.GetApiAuthenticationCookie(c))
		handler := getMapGroups.New(getMapGroups.NewModel(
			model.Name,
			model.ItemID,
			model.ProjectID,
		), a)

		return request.SendResponse[getMapGroups.Model](handler, c, http.StatusOK, func(c echo.Context, model interface{}) error {
			if a.ShouldRefresh() {
				session, err := a.Refresh()
				if err != nil {
					return err
				}

				c.SetCookie(request.EncryptAuthenticationCookie(session))
			}

			return nil
		}, false)
	}
}
