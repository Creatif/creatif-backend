package maps

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations/maps"
	"creatif/pkg/app/auth"
	addToMap2 "creatif/pkg/app/services/maps/addToMap"
	"creatif/pkg/lib/logger"
	"github.com/labstack/echo/v4"
	"net/http"
)

func AddToMapHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model maps.AddToMap
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = maps.SanitizeAddToMap(model)
		if model.Locale == "" {
			model.Locale = "eng"
		}

		l := logger.NewLogBuilder()
		handler := addToMap2.New(addToMap2.NewModel(model.ProjectID, model.Locale, model.Name, addToMap2.VariableModel{
			Name:      model.Entry.Name,
			Metadata:  []byte(model.Entry.Metadata),
			Groups:    model.Entry.Groups,
			Behaviour: model.Entry.Behaviour,
			Value:     []byte(model.Entry.Value),
		}), auth.NewApiAuthentication(), l)

		return request.SendResponse[addToMap2.Model](handler, c, http.StatusCreated, l)
	}
}
