package structures

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/app"
	"creatif/pkg/app/auth"
	removeStructure "creatif/pkg/app/services/structures/deleteStructure"
	"github.com/labstack/echo/v4"
	"net/http"
)

func RemoveStructureHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model app.RemoveStructure
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = app.SanitizeRemoveStructure(model)

		a := auth.NewApiAuthentication(request.GetApiAuthenticationCookie(c))
		handler := removeStructure.New(
			removeStructure.NewModel(model.ProjectID, model.ID, model.Type),
			a,
		)

		return request.SendResponse[removeStructure.Model](handler, c, http.StatusOK, func(c echo.Context, model interface{}) error {
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
