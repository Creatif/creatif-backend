package variables

import (
	declarations2 "creatif/cmd/http/handlers/declarations"
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations/variables"
	getVariable2 "creatif/pkg/app/services/variables/getVariable"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetVariableHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model variables.GetVariable
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = variables.SanitizeGetVariable(model)
		if model.Locale == "" {
			model.Locale = declarations2.DefaultLocale
		}

		handler := getVariable2.New(getVariable2.NewModel(model.ProjectID, model.Name, model.Locale, model.Fields))

		return request.SendResponse[getVariable2.Model](handler, c, http.StatusOK)
	}
}
