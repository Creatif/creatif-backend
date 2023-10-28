package lists

import (
	declarations2 "creatif/cmd/http/handlers/declarations"
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations/lists"
	"creatif/pkg/app/services/lists/deleteListItemByID"
	"github.com/labstack/echo/v4"
	"net/http"
)

func DeleteListItemByIDHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model lists.DeleteListItemByID
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = lists.SanitizeDeleteListItemByID(model)
		if model.Locale == "" {
			model.Locale = declarations2.DefaultLocale
		}

		handler := deleteListItemByID.New(deleteListItemByID.NewModel(
			model.ProjectID,
			model.Locale,
			model.Name,
			model.ItemID,
		))

		return request.SendResponse[deleteListItemByID.Model](handler, c, http.StatusOK)
	}
}
