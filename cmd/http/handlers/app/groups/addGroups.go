package groups

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/app"
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/groups/addGroups"
	"creatif/pkg/lib/logger"
	"creatif/pkg/lib/sdk"
	"github.com/labstack/echo/v4"
	"net/http"
)

func AddGroupsHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model app.AddGroups
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = app.SanitizeAddGroups(model)

		l := logger.NewLogBuilder()
		authentication := auth.NewApiAuthentication(request.GetApiAuthenticationCookie(c), l)
		handler := addGroups.New(addGroups.NewModel(model.ProjectID, sdk.Map(model.Groups, func(idx int, value app.SingleGroup) addGroups.GroupModel {
			return addGroups.GroupModel{
				ID:     value.ID,
				Name:   value.Name,
				Type:   value.Type,
				Action: value.Action,
			}
		})), authentication, l)

		return request.SendResponse[addGroups.Model](handler, c, http.StatusCreated, l, func(c echo.Context, model interface{}) error {
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
