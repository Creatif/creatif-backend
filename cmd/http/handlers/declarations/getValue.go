package declarations

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations"
	"creatif/pkg/app/services/getValue"
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

		handler := getValue.New(getValue.NewModel(model.ProjectID, model.Name))

		return request.SendResponse[getValue.Model](handler, c, http.StatusOK)
	}
}
