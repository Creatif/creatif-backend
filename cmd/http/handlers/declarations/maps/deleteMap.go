package maps

import (
	declarations2 "creatif/cmd/http/handlers/declarations"
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations"
	"creatif/pkg/app/services/maps/removeMap"
	"github.com/labstack/echo/v4"
	"net/http"
)

func DeleteMap() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model declarations.DeleteMap
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = declarations.SanitizeDeleteMap(model)
		if model.Locale == "" {
			model.Locale = declarations2.DefaultLocale
		}

		handler := removeMap.New(removeMap.NewModel(model.ProjectID, model.Locale, model.Name))

		return request.SendResponse[removeMap.Model](handler, c, http.StatusCreated)
	}
}
