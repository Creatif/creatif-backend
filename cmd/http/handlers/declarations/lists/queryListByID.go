package lists

import (
	declarations2 "creatif/cmd/http/handlers/declarations"
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations/lists"
	"creatif/pkg/app/services/lists/queryListByID"
	"creatif/pkg/lib/logger"
	"github.com/labstack/echo/v4"
	"net/http"
)

func QueryListByIDHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model lists.QueryListByID
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = lists.SanitizeQueryListByID(model)
		if model.Locale == "" {
			model.Locale = declarations2.DefaultLocale
		}

		l := logger.NewLogBuilder()
		handler := queryListByID.New(queryListByID.NewModel(
			model.ProjectID,
			model.Locale,
			model.Name,
			model.ID,
		), l)

		return request.SendResponse[queryListByID.Model](handler, c, http.StatusOK, l)
	}
}
