package declarations

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations"
	"creatif/pkg/app/services/updateMapVariable"
	"github.com/labstack/echo/v4"
	"net/http"
)

func UpdateMapVariableHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model declarations.UpdateMapVariable
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = declarations.SanitizeUpdateMapVariable(model)
		handler := updateMapVariable.New(updateMapVariable.NewModel(model.ProjectID, model.Name, updateMapVariable.VariableModel{
			Name:      model.Entry.Name,
			Metadata:  []byte(model.Entry.Metadata),
			Groups:    model.Entry.Groups,
			Behaviour: model.Entry.Behaviour,
			Value:     []byte(model.Entry.Value),
		}))

		return request.SendResponse[updateMapVariable.Model](handler, c, http.StatusOK)
	}
}
