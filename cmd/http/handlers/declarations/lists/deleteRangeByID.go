package lists

import (
	declarations2 "creatif/cmd/http/handlers/declarations"
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations/lists"
	"creatif/pkg/app/services/lists/deleteRangeByID"
	"github.com/labstack/echo/v4"
	"net/http"
)

func DeleteRangeByIDHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model lists.DeleteRangeByID
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = lists.SanitizeDeleteRangeByID(model)
		if model.Locale == "" {
			model.Locale = declarations2.DefaultLocale
		}

		handler := deleteRangeByID.New(deleteRangeByID.NewModel(
			model.ProjectID,
			model.Locale,
			model.Name,
			model.Items,
		))

		return request.SendResponse[deleteRangeByID.Model](handler, c, http.StatusCreated)
	}
}
