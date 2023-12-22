package variables

import (
	"creatif/cmd"
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations/variables"
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/variables/getVariableGroups"
	"creatif/pkg/lib/logger"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetVariableGroupsHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model variables.GetVariableGroups
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = variables.SanitizeGetVariableGroups(model)

		apiKey := c.Request().Header.Get(cmd.CreatifApiHeader)
		projectId := c.Request().Header.Get(cmd.CreatifProjectIDHeader)

		l := logger.NewLogBuilder()
		a := auth.NewApiAuthentication(request.GetApiAuthenticationCookie(c), projectId, apiKey, l)
		handler := getVariableGroups.New(getVariableGroups.NewModel(
			model.Name,
			model.ProjectID,
		), a, l)

		return request.SendResponse(handler, c, http.StatusOK, l, func(c echo.Context, model interface{}) error {
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
