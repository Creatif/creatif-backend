package maps

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations/maps"
	"creatif/pkg/app/auth"
	addToMap2 "creatif/pkg/app/services/maps/addToMap"
	"creatif/pkg/app/services/shared/connections"
	"creatif/pkg/lib/sdk"
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

		authentication := auth.NewApiAuthentication(request.GetApiAuthenticationCookie(c))
		handler := addToMap2.New(addToMap2.NewModel(model.ProjectID, model.Name, addToMap2.VariableModel{
			Name:      model.Variable.Name,
			Metadata:  []byte(model.Variable.Metadata),
			Locale:    model.Variable.Locale,
			Groups:    model.Variable.Groups,
			Behaviour: model.Variable.Behaviour,
			Value:     []byte(model.Variable.Value),
		}, sdk.Map(model.Connections, func(idx int, value maps.Connection) connections.Connection {
			return connections.Connection{
				Path:          value.Path,
				StructureType: value.StructureType,
				VariableID:    value.VariableID,
			}
		}), model.ImagePaths), authentication)

		return request.SendResponse[addToMap2.Model](handler, c, http.StatusCreated, func(c echo.Context, model interface{}) error {
			if authentication.ShouldRefresh() {
				session, err := authentication.Refresh()
				if err != nil {
					return err
				}

				c.SetCookie(request.EncryptAuthenticationCookie(session))
			}

			return nil
		}, false)
	}
}
