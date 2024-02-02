package lists

import (
	"creatif/cmd"
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations/lists"
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/lists/updateListItemByID"
	"creatif/pkg/app/services/shared"
	"creatif/pkg/lib/logger"
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

		apiKey := c.Request().Header.Get(cmd.CreatifApiHeader)
		projectId := c.Request().Header.Get(cmd.CreatifProjectIDHeader)

		l := logger.NewLogBuilder()
		authentication := auth.NewApiAuthentication(request.GetApiAuthenticationCookie(c), projectId, apiKey, l)
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
		), authentication, l)

		res := request.SendResponse[updateListItemByID.Model](handler, c, http.StatusOK, l, func(c echo.Context, model interface{}) error {
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
