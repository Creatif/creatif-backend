package maps

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations/maps"
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/maps/removeMapVariable"
	"github.com/labstack/echo/v4"
	"net/http"
)

func DeleteMapEntry() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model maps.DeleteMapEntry
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = maps.SanitizeDeleteMapEntry(model)

		authentication := auth.NewApiAuthentication(request.GetApiAuthenticationCookie(c))
		handler := removeMapVariable.New(removeMapVariable.NewModel(model.ProjectID, model.Name, model.VariableName), authentication)

		return request.SendResponse[removeMapVariable.Model](handler, c, http.StatusOK, func(c echo.Context, model interface{}) error {
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
