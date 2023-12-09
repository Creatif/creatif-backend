package lists

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations/lists"
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/lists/switchByID"
	"creatif/pkg/lib/logger"
	"github.com/labstack/echo/v4"
	"net/http"
)

func SwitchByIDHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model lists.SwitchByID
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = lists.SanitizeSwitchByID(model)

		l := logger.NewLogBuilder()
		handler := switchByID.New(switchByID.NewModel(
			model.ProjectID,
			model.Name,
			model.ID,
			model.ShortID,
			model.Source,
			model.Destination,
		), auth.NewNoopAuthentication(), l)

		return request.SendResponse[switchByID.Model](handler, c, http.StatusOK, l, nil, false)
	}
}
