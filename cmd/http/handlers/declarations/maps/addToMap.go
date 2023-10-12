package maps

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations"
	addToMap2 "creatif/pkg/app/services/maps/addToMap"
	"github.com/labstack/echo/v4"
	"net/http"
)

func AddToMapHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model declarations.AddToMap
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = declarations.SanitizeAddToMap(model)
		handler := addToMap2.New(addToMap2.NewModel(model.ProjectID, model.Name, addToMap2.VariableModel{
			Name:      model.Entry.Name,
			Metadata:  []byte(model.Entry.Metadata),
			Groups:    model.Entry.Groups,
			Behaviour: model.Entry.Behaviour,
			Value:     []byte(model.Entry.Value),
		}))

		return request.SendResponse[addToMap2.Model](handler, c, http.StatusCreated)
	}
}
