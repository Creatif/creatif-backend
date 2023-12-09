package lists

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations/lists"
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/lists/updateList"
	"creatif/pkg/lib/logger"
	"github.com/labstack/echo/v4"
	"net/http"
)

func UpdateListHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model lists.UpdateList
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = lists.SanitizeUpdateList(model)

		l := logger.NewLogBuilder()
		handler := updateList.New(updateList.NewModel(
			model.ProjectID,
			model.Fields,
			model.Name,
			model.ID,
			model.ShortID,
			model.Values.Name,
		), auth.NewNoopAuthentication(), l)

		return request.SendResponse[updateList.Model](handler, c, http.StatusOK, l, nil, false)
	}
}
