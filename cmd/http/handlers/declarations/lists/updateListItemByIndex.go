package lists

import (
	declarations2 "creatif/cmd/http/handlers/declarations"
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations/lists"
	"creatif/pkg/app/services/lists/updateListItemByIndex"
	"github.com/labstack/echo/v4"
	"net/http"
)

func UpdateListItemByIndexHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model lists.UpdateListItemByIndex
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = lists.SanitizeUpdateListItemByIndex(model)
		if model.Locale == "" {
			model.Locale = declarations2.DefaultLocale
		}

		handler := updateListItemByIndex.New(updateListItemByIndex.NewModel(
			model.ProjectID,
			model.Locale,
			model.Fields,
			model.ListName,
			model.Index,
			model.Values.Name,
			model.Values.Behaviour,
			model.Values.Groups,
			[]byte(model.Values.Metadata),
			[]byte(model.Values.Value),
		))

		return request.SendResponse[updateListItemByIndex.Model](handler, c, http.StatusCreated)
	}
}
