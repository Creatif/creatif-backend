package declarations

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations"
	create "creatif/pkg/app/declarations/createVariable"
	"github.com/labstack/echo/v4"
	"net/http"
)

func CreateVariableHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model declarations.CreateVariable
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = declarations.SanitizeVariable(model)

		handler := create.New(create.NewModel(model.Name, model.Behaviour, model.Groups, []byte(model.Metadata), []byte(model.Value)))

		return request.SendResponse[create.Model](handler, c, http.StatusCreated)
	}
}
