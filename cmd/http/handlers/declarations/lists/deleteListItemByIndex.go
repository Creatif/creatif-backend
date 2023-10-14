package lists

import (
	declarations2 "creatif/cmd/http/handlers/declarations"
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations/lists"
	"creatif/pkg/app/services/lists/deleteListItemByIndex"
	"github.com/labstack/echo/v4"
	"net/http"
)

func DeleteListItemByIndexHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model lists.DeleteListItemByIndex
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = lists.SanitizeDeleteListItemByIndex(model)
		if model.Locale == "" {
			model.Locale = declarations2.DefaultLocale
		}

		handler := deleteListItemByIndex.New(deleteListItemByIndex.NewModel(
			model.ProjectID,
			model.Locale,
			model.Name,
			model.ItemIndex,
		))

		return request.SendResponse[deleteListItemByIndex.Model](handler, c, http.StatusCreated)
	}
}
