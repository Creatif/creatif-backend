package variables

import (
	declarations2 "creatif/cmd/http/handlers/declarations"
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations/variables"
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/variables/deleteVariable"
	"creatif/pkg/lib/logger"
	"github.com/labstack/echo/v4"
	"net/http"
)

func DeleteVariableHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model variables.DeleteVariable
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = variables.SanitizeDeleteVariable(model)
		if model.Locale == "" {
			model.Locale = declarations2.DefaultLocale
		}

		l := logger.NewLogBuilder()
		handler := deleteVariable.New(deleteVariable.NewModel(model.ProjectID, model.ID, model.ShortID, model.Name, model.Locale), auth.NewNoopAuthentication(), l)

		return request.SendResponse[deleteVariable.Model](handler, c, http.StatusOK, l, nil)
	}
}
