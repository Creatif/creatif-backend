package maps

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations/maps"
	"creatif/pkg/app/auth"
	updateMapVariable2 "creatif/pkg/app/services/maps/updateMapVariable"
	"creatif/pkg/app/services/shared"
	"creatif/pkg/lib/logger"
	"creatif/pkg/lib/sdk"
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
		references := make([]shared.UpdateReference, 0)
		if len(model.References) > 0 {
			references = sdk.Map(model.References, func(idx int, value maps.UpdateReference) shared.UpdateReference {
				return shared.UpdateReference{
					Name:          value.Name,
					StructureName: value.StructureName,
					StructureType: value.StructureType,
					VariableID:    value.VariableID,
				}
			})
		}

		l := logger.NewLogBuilder()
		handler := updateMapVariable2.New(updateMapVariable2.NewModel(model.ProjectID, model.Name, model.ItemID, model.ResolvedFields, updateMapVariable2.VariableModel{
			Name:      model.Variable.Name,
			Metadata:  []byte(model.Variable.Metadata),
			Locale:    model.Variable.Locale,
			Groups:    model.Variable.Groups,
			Behaviour: model.Variable.Behaviour,
			Value:     []byte(model.Variable.Value),
		}, references), auth.NewNoopAuthentication(), l)

		return request.SendResponse[updateMapVariable2.Model](handler, c, http.StatusOK, l, nil, false)
	}
}
