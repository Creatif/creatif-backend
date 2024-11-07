package maps

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations/maps"
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/maps/switchByID"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

func SwitchByIDHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model maps.SwitchByID
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = maps.SanitizeSwitchByID(model)

		authentication := auth.NewApiAuthentication(request.GetApiAuthenticationCookie(c))
		handler := switchByID.New(switchByID.NewModel(
			model.ProjectID,
			model.Name,
			model.Source,
			model.Destination,
			strings.ToLower(model.OrderDirection),
		), authentication)

		return request.SendResponse[switchByID.Model](handler, c, http.StatusOK, func(c echo.Context, model interface{}) error {
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
