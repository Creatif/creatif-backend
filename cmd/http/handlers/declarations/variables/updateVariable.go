package variables

import (
	"creatif/cmd"
	declarations2 "creatif/cmd/http/handlers/declarations"
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations/variables"
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/variables/updateVariable"
	"creatif/pkg/lib/logger"
	"github.com/labstack/echo/v4"
	"net/http"
)

func UpdateVariableHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model variables.UpdateVariable
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = variables.SanitizeUpdateVariable(model)
		if model.Locale == "" {
			model.Locale = declarations2.DefaultLocale
		}

		apiKey := c.Request().Header.Get(cmd.CreatifApiHeader)
		projectId := c.Request().Header.Get(cmd.CreatifProjectIDHeader)

		l := logger.NewLogBuilder()
		authentication := auth.NewApiAuthentication(request.GetApiAuthenticationCookie(c), projectId, apiKey, l)
		handler := updateVariable.New(updateVariable.NewModel(
			model.ProjectID,
			model.Fields,
			model.Name,
			model.Values.Name,
			model.Values.Behaviour,
			model.Values.Groups,
			[]byte(model.Values.Metadata),
			[]byte(model.Values.Value),
			model.Values.Locale,
		), authentication, l)

		return request.SendResponse[updateVariable.Model](handler, c, http.StatusOK, l, func(c echo.Context, model interface{}) error {
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
