package groups

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/app"
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/groups/getGroups"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetGroupsHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model app.GetGroups
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = app.SanitizeGetGroups(model)

		authentication := auth.NewApiAuthentication(request.GetApiAuthenticationCookie(c))
		handler := getGroups.New(getGroups.NewModel(model.ProjectID), authentication)

		return request.SendResponse[getGroups.Model](handler, c, http.StatusOK, func(c echo.Context, model interface{}) error {
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
