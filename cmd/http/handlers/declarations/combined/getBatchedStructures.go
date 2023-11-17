package combined

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations/combined"
	"creatif/pkg/app/auth"
	getBatchStructures2 "creatif/pkg/app/services/combined/getBatchStructures"
	"creatif/pkg/lib/logger"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetBatchedStructuresHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model combined.GetBatchedStructures
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = combined.SanitizeGetBatchedVariables(model)
		serviceModel := make(map[string]string)
		for _, m := range model.Structures {
			serviceModel[m.Name] = m.Type
		}

		l := logger.NewLogBuilder()
		handler := getBatchStructures2.New(getBatchStructures2.NewModel(model.ProjectID, serviceModel), auth.NewNoopAuthentication(), l)

		return request.SendResponse[*getBatchStructures2.Model](handler, c, http.StatusCreated, l, nil)
	}
}
