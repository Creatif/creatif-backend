package declarations

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations"
	"creatif/pkg/app/combined/getBatchData"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetBatchedVariablesHandler() func(e echo.Context) error {
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

		handler := getBatchData.New(getBatchData.NewModel(serviceModel))

		return request.SendResponse[*getBatchData.Model](handler, c, http.StatusCreated)
	}
}
