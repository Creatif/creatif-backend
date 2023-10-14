package lists

import (
	declarations2 "creatif/cmd/http/handlers/declarations"
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations/lists"
	"creatif/pkg/app/services/lists/switchByIndex"
	"github.com/labstack/echo/v4"
	"net/http"
)

func SwitchByIndexHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model lists.SwitchByIndex
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = lists.SanitizeSwitchByIndex(model)
		if model.Locale == "" {
			model.Locale = declarations2.DefaultLocale
		}

		handler := switchByIndex.New(switchByIndex.NewModel(
			model.ProjectID,
			model.Locale,
			model.Name,
			model.Source,
			model.Destination,
		))

		return request.SendResponse[switchByIndex.Model](handler, c, http.StatusCreated)
	}
}
