package declarations

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations"
	getVariable2 "creatif/pkg/app/services/variables/getVariable"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetVariableHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model declarations.GetVariable
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = declarations.SanitizeGetVariable(model)

		handler := getVariable2.New(getVariable2.NewModel(model.ProjectID, model.Name, model.Fields))

		return request.SendResponse[getVariable2.Model](handler, c, http.StatusCreated)
	}
}
