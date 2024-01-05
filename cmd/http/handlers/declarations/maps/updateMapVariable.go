package maps

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations/maps"
	"creatif/pkg/app/auth"
	updateMapVariable2 "creatif/pkg/app/services/maps/updateMapVariable"
	"creatif/pkg/lib/logger"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

func UpdateMapVariableHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model maps.UpdateMapVariable
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		model.Fields = c.QueryParam("fields")

		model = maps.SanitizeUpdateMapVariable(model)

		fmt.Println(model.ResolvedFields)
		l := logger.NewLogBuilder()
		handler := updateMapVariable2.New(updateMapVariable2.NewModel(model.ProjectID, model.Name, model.ItemID, model.ResolvedFields, updateMapVariable2.VariableModel{
			Name:      model.Variable.Name,
			Metadata:  []byte(model.Variable.Metadata),
			Locale:    model.Variable.Locale,
			Groups:    model.Variable.Groups,
			Behaviour: model.Variable.Behaviour,
			Value:     []byte(model.Variable.Value),
		}), auth.NewNoopAuthentication(), l)

		return request.SendResponse[updateMapVariable2.Model](handler, c, http.StatusOK, l, nil, false)
	}
}
