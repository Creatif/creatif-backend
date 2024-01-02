package maps

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations/maps"
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/maps/removeMapVariable"
	"creatif/pkg/lib/logger"
	"github.com/labstack/echo/v4"
	"net/http"
)

func DeleteMapEntry() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model maps.DeleteMapEntry
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = maps.SanitizeDeleteMapEntry(model)

		l := logger.NewLogBuilder()
		handler := removeMapVariable.New(removeMapVariable.NewModel(model.ProjectID, model.Name, model.VariableName), auth.NewNoopAuthentication(), l)

		return request.SendResponse[removeMapVariable.Model](handler, c, http.StatusOK, l, nil, false)
	}
}
