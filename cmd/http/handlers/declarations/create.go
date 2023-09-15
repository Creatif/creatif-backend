package declarations

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations"
	create "creatif/pkg/app/declarations/createNode"
	"github.com/labstack/echo/v4"
	"net/http"
)

func CreateNodeHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model declarations.CreateNode
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = declarations.SanitizeNode(model)

		handler := create.New(create.NewModel(model.Name, model.Behaviour, model.Groups, []byte(model.Metadata)))

		return request.SendResponse[create.Model](handler, c, http.StatusCreated)
	}
}
