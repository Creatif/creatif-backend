package declarations

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations"
	"creatif/pkg/app/declarations/getMap"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

func GetMapHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model declarations.GetMap
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = declarations.SanitizeGetMap(model)

		fields := strings.Split(model.Fields, ",")
		newFields := make([]string, 0)
		for _, f := range fields {
			newFields = append(newFields, strings.Trim(f, " "))
		}

		handler := getMap.New(getMap.NewModel(model.ID, newFields))

		return request.SendResponse[getMap.Model](handler, c, http.StatusCreated)
	}
}
