package declarations

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations"
	"creatif/pkg/app/combined/getBatchNodes"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetBatchedNodesHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model []declarations.GetBatchedNodes
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = declarations.SanitizeGetBatchedNodes(model)
		serviceModel := make(map[string]string)
		for _, m := range model {
			serviceModel[m.Name] = m.Type
		}

		handler := getBatchNodes.New(getBatchNodes.NewModel(serviceModel))

		return request.SendResponse[*getBatchNodes.Model](handler, c, http.StatusCreated)
	}
}
