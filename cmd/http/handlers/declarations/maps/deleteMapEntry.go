package maps

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations/maps"
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/maps/removeMapEntry"
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
		if model.Locale == "" {
			model.Locale = "eng"
		}

		l := logger.NewLogBuilder()
		handler := removeMapEntry.New(removeMapEntry.NewModel(model.ProjectID, model.Locale, model.Name, model.MapID, model.MapShortID, model.VariableName, model.VariableID, model.VariableShortID), auth.NewNoopAuthentication(), l)

		return request.SendResponse[removeMapEntry.Model](handler, c, http.StatusOK, l, nil, false)
	}
}
