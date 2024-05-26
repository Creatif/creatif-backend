package getMany

import (
	publicApi2 "creatif/cmd/http/handlers/publicApi"
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/publicApi"
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/publicApi/getManyItems"
	"creatif/pkg/lib/logger"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetManyHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model publicApi.GetMany
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		versionName := c.Request().Header.Get(publicApi2.CreatifVersionHeader)
		model.VersionName = versionName
		model = publicApi.SanitizeGetMany(model)

		l := logger.NewLogBuilder()
		handler := getManyItems.New(getManyItems.NewModel(model.VersionName, model.ProjectID, model.ResolvedIds, getManyItems.Options{
			ValueOnly: model.ResolvedOptions.ValueOnly,
		}), auth.NewAnonymousAuthentication(), l)

		return request.SendPublicResponse[getManyItems.Model](handler, c, http.StatusOK, l, nil, false)
	}
}
