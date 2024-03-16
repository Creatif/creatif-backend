package getListItemByName

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/publicApi"
	"creatif/pkg/app/auth"
	getListItemsByNameService "creatif/pkg/app/services/publicApi/getListItemsByName"
	"creatif/pkg/lib/logger"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetListItemByNameHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model publicApi.GetListItemByName
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = publicApi.SanitizeGetListItemByName(model)

		l := logger.NewLogBuilder()
		handler := getListItemsByNameService.New(getListItemsByNameService.NewModel(model.ProjectID, model.StructureName, model.Name, model.Locale, getListItemsByNameService.Options{
			ValueOnly: model.ResolvedOptions.ValueOnly,
		}), auth.NewAnonymousAuthentication(), l)

		return request.SendPublicResponse[getListItemsByNameService.Model](handler, c, http.StatusOK, l, nil, false)
	}
}
