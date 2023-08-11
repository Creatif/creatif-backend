package handlers

import (
	"creatif/cmd/http/request"
	create "creatif/pkg/app/nodes/create"
	"github.com/labstack/echo/v4"
	"net/http"
)

func CreateNodeHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model request.CreateNode
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		model = request.SanitizeNode(model)

		handler := create.New(create.NewCreateNodeModel(model.Name))

		return request.SendResponse[create.CreateNodeModel](handler, c, http.StatusCreated)
	}
}
