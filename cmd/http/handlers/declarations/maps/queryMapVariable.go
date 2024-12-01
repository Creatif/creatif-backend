package maps

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations/maps"
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/maps/queryMapVariable"
	"github.com/labstack/echo/v4"
	"net/http"
)

func QueryMapVariableHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model maps.QueryMapVariable
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = maps.SanitizeQueryMapVariable(model)

		authentication := auth.NewApiAuthentication(request.GetApiAuthenticationCookie(c))
		handler := queryMapVariable.New(queryMapVariable.NewModel(
			model.ProjectID,
			model.Name,
			model.ItemID,
			model.ConnectionViewType,
		), authentication)

		return request.SendResponse[queryMapVariable.Model](handler, c, http.StatusOK, func(c echo.Context, model interface{}) error {
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
