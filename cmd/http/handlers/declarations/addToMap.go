package declarations

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations"
	"creatif/pkg/app/services/addToMap"
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
		handler := addToMap.New(addToMap.NewModel(model.ProjectID, model.Name, addToMap.VariableModel{
			Name:      model.Entry.Name,
			Metadata:  []byte(model.Entry.Metadata),
			Groups:    model.Entry.Groups,
			Behaviour: model.Entry.Behaviour,
			Value:     []byte(model.Entry.Value),
		}))

		return request.SendResponse[addToMap.Model](handler, c, http.StatusCreated)
	}
}
