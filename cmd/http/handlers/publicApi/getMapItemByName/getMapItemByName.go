package getMapItemByName

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/publicApi"
	"creatif/pkg/app/auth"
	getMapItemByNameService "creatif/pkg/app/services/publicApi/getMapItemByName"
	"creatif/pkg/lib/logger"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetMapItemByNameHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model publicApi.GetMapItemByName
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = publicApi.SanitizeGetMapItemByName(model)

		l := logger.NewLogBuilder()
		handler := getMapItemByNameService.New(getMapItemByNameService.NewModel(model.ProjectID, model.StructureName, model.Name, model.Locale, getMapItemByNameService.Options{
			ValueOnly: model.ResolvedOptions.ValueOnly,
		}), auth.NewAnonymousAuthentication(), l)

		return request.SendPublicResponse[getMapItemByNameService.Model](handler, c, http.StatusOK, l, nil, false)
	}
}
