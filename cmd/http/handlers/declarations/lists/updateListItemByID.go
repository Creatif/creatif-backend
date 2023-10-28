package lists

import (
	declarations2 "creatif/cmd/http/handlers/declarations"
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations/lists"
	"creatif/pkg/app/services/lists/updateListItemByID"
	"github.com/labstack/echo/v4"
	"net/http"
)

func UpdateListItemByIDHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model lists.UpdateListItemByID
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = lists.SanitizeUpdateListItemByID(model)
		if model.Locale == "" {
			model.Locale = declarations2.DefaultLocale
		}

		handler := updateListItemByID.New(updateListItemByID.NewModel(
			model.ProjectID,
			model.Locale,
			model.Fields,
			model.ListName,
			model.ItemID,
			model.Values.Name,
			model.Values.Behaviour,
			model.Values.Groups,
			[]byte(model.Values.Metadata),
			[]byte(model.Values.Value),
		))

		return request.SendResponse[updateListItemByID.Model](handler, c, http.StatusOK)
	}
}
