package maps

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations/maps"
	"creatif/pkg/app/auth"
	updateMapVariable2 "creatif/pkg/app/services/maps/updateMapVariable"
	"creatif/pkg/lib/logger"
	"github.com/labstack/echo/v4"
	"net/http"
)

func UpdateMapVariableHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model maps.UpdateMapVariable
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = maps.SanitizeUpdateMapVariable(model)
		if model.Locale == "" {
			model.Locale = "eng"
		}

		l := logger.NewLogBuilder()
		handler := updateMapVariable2.New(updateMapVariable2.NewModel(model.ProjectID, model.Locale, model.MapName, model.VariableName, model.SanitizedFields, updateMapVariable2.VariableModel{
			Name:      model.Entry.Name,
			Metadata:  []byte(model.Entry.Metadata),
			Groups:    model.Entry.Groups,
			Behaviour: model.Entry.Behaviour,
			Value:     []byte(model.Entry.Value),
		}), auth.NewApiAuthentication(), l)

		return request.SendResponse[updateMapVariable2.Model](handler, c, http.StatusOK, l)
	}
}
