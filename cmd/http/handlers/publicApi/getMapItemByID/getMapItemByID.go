package getMapItemByID

import (
	publicApi2 "creatif/cmd/http/handlers/publicApi"
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/publicApi"
	"creatif/pkg/app/auth"
	getMapItemByIDService "creatif/pkg/app/services/publicApi/getMapItemById"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetMapItemByIDHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model publicApi.GetMapItemByID
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		versionName := c.Request().Header.Get(publicApi2.CreatifVersionHeader)
		model.VersionName = versionName
		model = publicApi.SanitizeGetMapItemByID(model)

		handler := getMapItemByIDService.New(getMapItemByIDService.NewModel(model.VersionName, model.ProjectID, model.ItemID, getMapItemByIDService.Options{
			ValueOnly: model.ResolvedOptions.ValueOnly,
		}), auth.NewAnonymousAuthentication())

		return request.SendPublicResponse[getMapItemByIDService.Model](handler, c, http.StatusOK, nil, false)
	}
}
