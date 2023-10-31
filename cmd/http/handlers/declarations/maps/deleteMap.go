package maps

import (
	declarations2 "creatif/cmd/http/handlers/declarations"
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations/maps"
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/maps/removeMap"
	"creatif/pkg/lib/logger"
	"github.com/labstack/echo/v4"
	"net/http"
)

func DeleteMap() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model maps.DeleteMap
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = maps.SanitizeDeleteMap(model)
		if model.Locale == "" {
			model.Locale = declarations2.DefaultLocale
		}

		l := logger.NewLogBuilder()
		handler := removeMap.New(removeMap.NewModel(model.ProjectID, model.Locale, model.Name, model.ID, model.ShortID), auth.NewApiAuthentication(), l)

		return request.SendResponse[removeMap.Model](handler, c, http.StatusOK, l)
	}
}
