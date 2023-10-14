package lists

import (
	declarations2 "creatif/cmd/http/handlers/declarations"
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations/lists"
	"creatif/pkg/app/services/lists/queryListByIndex"
	"github.com/labstack/echo/v4"
	"net/http"
)

func QueryListByIndexHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model lists.QueryListByIndex
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = lists.SanitizeQueryListByIndex(model)
		if model.Locale == "" {
			model.Locale = declarations2.DefaultLocale
		}

		handler := queryListByIndex.New(queryListByIndex.NewModel(
			model.ProjectID,
			model.Locale,
			model.Name,
			model.Index,
		))

		return request.SendResponse[queryListByIndex.Model](handler, c, http.StatusCreated)
	}
}
