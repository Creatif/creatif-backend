package declarations

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations"
	updateVariable "creatif/pkg/app/declarations/updateVariable"
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

		handler := updateVariable.New(updateVariable.NewModel(model.Fields, model.Name, model.Values.Name, model.Values.Behaviour, model.Values.Groups, model.Values.Metadata, model.Values.Value))

		return request.SendResponse[updateVariable.Model](handler, c, http.StatusOK)
	}
}
