package declarations

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations"
	"creatif/pkg/app/combined/getBatchStructures"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetBatchedStructuresHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model []declarations.GetBatchedVariables
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = declarations.SanitizeGetBatchedVariables(model)
		serviceModel := make(map[string]string)
		for _, m := range model {
			serviceModel[m.Name] = m.Type
		}

		handler := getBatchStructures.New(getBatchStructures.NewModel(serviceModel))

		return request.SendResponse[*getBatchStructures.Model](handler, c, http.StatusCreated)
	}
}
