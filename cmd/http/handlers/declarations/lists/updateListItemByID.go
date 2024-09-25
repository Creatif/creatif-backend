package lists

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations/lists"
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/lists/updateListItemByID"
	"creatif/pkg/app/services/shared"
	"creatif/pkg/lib/sdk"
	"github.com/labstack/echo/v4"
	"net/http"
)

func UpdateListItemByIDHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model lists.UpdateListItemByID
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		// hack! query tag just won't work
		model.Fields = c.QueryParam("fields")

		model = lists.SanitizeUpdateListItemByID(model)
		references := make([]shared.UpdateReference, 0)
		if len(model.References) > 0 {
			references = sdk.Map(model.References, func(idx int, value lists.UpdateReference) shared.UpdateReference {
				return shared.UpdateReference{
					Name:          value.Name,
					StructureName: value.StructureName,
					StructureType: value.StructureType,
					VariableID:    value.VariableID,
				}
			})
		}

		authentication := auth.NewApiAuthentication(request.GetApiAuthenticationCookie(c))
		handler := updateListItemByID.New(updateListItemByID.NewModel(
			model.ProjectID,
			model.Values.Locale,
			model.ResolvedFields,
			model.Name,
			model.ItemID,
			model.Values.Name,
			model.Values.Behaviour,
			model.Values.Groups,
			[]byte(model.Values.Metadata),
			[]byte(model.Values.Value),
			references,
			model.ImagePaths,
		), authentication)

		res := request.SendResponse[updateListItemByID.Model](handler, c, http.StatusOK, func(c echo.Context, model interface{}) error {
			if authentication.ShouldRefresh() {
				session, err := authentication.Refresh()
				if err != nil {
					return err
				}

				c.SetCookie(request.EncryptAuthenticationCookie(session))
			}

			return nil
		}, false)

		return res
	}
}
