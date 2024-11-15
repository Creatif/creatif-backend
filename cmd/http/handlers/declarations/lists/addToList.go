package lists

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations/lists"
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/lists/addToList"
	"creatif/pkg/app/services/shared/connections"
	"creatif/pkg/lib/sdk"
	"github.com/labstack/echo/v4"
	"net/http"
)

func AddToListHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model lists.AddToList
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = lists.SanitizeAddToList(model)

		authentication := auth.NewApiAuthentication(request.GetApiAuthenticationCookie(c))
		handler := addToList.New(addToList.NewModel(model.ProjectID, model.Name, addToList.VariableModel{
			Name:      model.Variable.Name,
			Metadata:  []byte(model.Variable.Metadata),
			Locale:    model.Variable.Locale,
			Groups:    model.Variable.Groups,
			Behaviour: model.Variable.Behaviour,
			Value:     []byte(model.Variable.Value),
		}, sdk.Map(model.Connections, func(idx int, value lists.Connection) connections.Connection {
			return connections.Connection{
				Path:          value.Name,
				StructureType: value.StructureType,
				VariableID:    value.VariableID,
			}
		}), model.ImagePaths), authentication)

		return request.SendResponse[addToList.Model](handler, c, http.StatusCreated, func(c echo.Context, model interface{}) error {
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
