package declarations

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations"
	delete "creatif/pkg/app/declarations/deleteVariable"
	"github.com/labstack/echo/v4"
	"net/http"
)

func DeleteVariableHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model declarations.DeleteVariable
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = declarations.SanitizeDeleteVariable(model)

		handler := delete.New(delete.NewModel(model.Name))

		return request.SendResponse[delete.Model](handler, c, http.StatusCreated)
	}
}
