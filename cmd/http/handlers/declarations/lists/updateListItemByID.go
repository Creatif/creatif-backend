package lists

import (
	declarations2 "creatif/cmd/http/handlers/declarations"
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations/lists"
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/lists/updateListItemByID"
	"creatif/pkg/lib/logger"
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

		l := logger.NewLogBuilder()
		handler := updateListItemByID.New(updateListItemByID.NewModel(
			model.ProjectID,
			model.Locale,
			model.Fields,
			model.ListName,
			model.ListID,
			model.ListShortID,
			model.ItemID,
			model.ItemShortID,
			model.Values.Name,
			model.Values.Behaviour,
			model.Values.Groups,
			[]byte(model.Values.Metadata),
			[]byte(model.Values.Value),
		), auth.NewNoopAuthentication(), l)

		return request.SendResponse[updateListItemByID.Model](handler, c, http.StatusOK, l, nil)
	}
}
