package declarations

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations"
	"creatif/pkg/app/declarations/getNode"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetNodeHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model declarations.GetNode
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = declarations.SanitizeGetNode(model)

		handler := getNode.New(getNode.NewGetNodeModel(model.ID, model.Fields))

		return request.SendResponse[getNode.GetNodeModel](handler, c, http.StatusCreated)
	}
}
