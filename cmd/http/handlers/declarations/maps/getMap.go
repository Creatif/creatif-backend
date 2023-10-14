package maps

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations/maps"
	getMap2 "creatif/pkg/app/services/maps/getMap"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

func GetMapHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model maps.GetMap
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = maps.SanitizeGetMap(model)
		if model.Locale == "" {
			model.Locale = "eng"
		}

		newFields := make([]string, 0)
		if strings.Trim(model.Fields, " ") != "" {
			fields := strings.Split(strings.Trim(model.Fields, " "), ",")
			for _, f := range fields {
				newFields = append(newFields, strings.Trim(f, " "))
			}
		}

		handler := getMap2.New(getMap2.NewModel(model.ProjectID, model.Locale, model.Name, newFields))

		return request.SendResponse[getMap2.Model](handler, c, http.StatusOK)
	}
}
