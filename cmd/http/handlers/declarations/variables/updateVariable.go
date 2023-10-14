package variables

import (
	declarations2 "creatif/cmd/http/handlers/declarations"
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations"
	"creatif/pkg/app/services/variables/updateVariable"
	"github.com/labstack/echo/v4"
	"net/http"
)

func UpdateVariableHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model declarations.UpdateVariable
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = declarations.SanitizeUpdateVariable(model)
		if model.Locale == "" {
			model.Locale = declarations2.DefaultLocale
		}

		handler := updateVariable.New(updateVariable.NewModel(
			model.ProjectID,
			model.Locale,
			model.Fields,
			model.Name,
			model.Values.Name,
			model.Values.Behaviour,
			model.Values.Groups,
			[]byte(model.Values.Metadata),
			[]byte(model.Values.Value)),
		)

		return request.SendResponse[updateVariable.Model](handler, c, http.StatusOK)
	}
}
