package variables

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations"
	getValue2 "creatif/pkg/app/services/variables/getValue"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetValueHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model declarations.GetValue
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = declarations.SanitizeGetValue(model)
		if model.Locale == "" {
			model.Locale = "eng"
		}

		handler := getValue2.New(getValue2.NewModel(model.ProjectID, model.Name, model.Locale))

		return request.SendResponse[getValue2.Model](handler, c, http.StatusOK)
	}
}
