package variables

import (
	declarations2 "creatif/cmd/http/handlers/declarations"
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations/variables"
	"creatif/pkg/app/auth"
	getValue2 "creatif/pkg/app/services/variables/getValue"
	"creatif/pkg/lib/logger"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetValueHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model variables.GetValue
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = variables.SanitizeGetValue(model)
		if model.Locale == "" {
			model.Locale = declarations2.DefaultLocale
		}

		l := logger.NewLogBuilder()
		handler := getValue2.New(getValue2.NewModel(model.ProjectID, model.Name, model.Locale), auth.NewApiAuthentication(), l)

		return request.SendResponse[getValue2.Model](handler, c, http.StatusOK, l)
	}
}
