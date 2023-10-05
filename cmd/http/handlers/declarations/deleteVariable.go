package declarations

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations"
	"creatif/pkg/app/services/variables/deleteVariable"
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

		handler := deleteVariable.New(deleteVariable.NewModel(model.ProjectID, model.Name))

		return request.SendResponse[deleteVariable.Model](handler, c, http.StatusCreated)
	}
}
