package variables

import (
	declarations2 "creatif/cmd/http/handlers/declarations"
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations/variables"
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/variables/createVariable"
	"creatif/pkg/lib/logger"
	"github.com/labstack/echo/v4"
	"net/http"
)

func CreateVariableHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model variables.CreateVariable
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = variables.SanitizeVariable(model)
		if model.Locale == "" {
			model.Locale = declarations2.DefaultLocale
		}

		l := logger.NewLogBuilder()
		handler := createVariable.New(createVariable.NewModel(model.ProjectID, model.Locale, model.Name, model.Behaviour, model.Groups, []byte(model.Metadata), []byte(model.Value)), auth.NewApiAuthentication(), l)

		return request.SendResponse[createVariable.Model](handler, c, http.StatusCreated, l, nil)
	}
}
