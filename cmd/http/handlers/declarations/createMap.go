package declarations

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations"
	"creatif/pkg/app/declarations/mapCreate"
	"github.com/labstack/echo/v4"
	"net/http"
)

func CreateMapHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model declarations.CreateMap
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = declarations.SanitizeMapModel(model)

		handler := mapCreate.New(mapCreate.NewModel(model.Name, model.Nodes))

		return request.SendResponse[mapCreate.Model](handler, c, http.StatusCreated)
	}
}
