package maps

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations/maps"
	"creatif/pkg/app/auth"
	getMap2 "creatif/pkg/app/services/maps/getMap"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetMapHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model maps.GetMap
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = maps.SanitizeGetMap(model)

		authentication := auth.NewApiAuthentication(request.GetApiAuthenticationCookie(c))
		handler := getMap2.New(getMap2.NewModel(model.ProjectID, model.Name), authentication)

		return request.SendResponse[getMap2.Model](handler, c, http.StatusOK, func(c echo.Context, model interface{}) error {
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
