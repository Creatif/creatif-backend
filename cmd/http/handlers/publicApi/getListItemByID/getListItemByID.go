package getListItemByID

import (
	publicApi2 "creatif/cmd/http/handlers/publicApi"
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/publicApi"
	"creatif/pkg/app/auth"
	getListItemByIDService "creatif/pkg/app/services/publicApi/getListItemById"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetListItemByIDHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model publicApi.GetListItemByID
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		versionName := c.Request().Header.Get(publicApi2.CreatifVersionHeader)
		model.VersionName = versionName

		model = publicApi.SanitizeGetListItemByID(model)

		handler := getListItemByIDService.New(getListItemByIDService.NewModel(model.VersionName, model.ProjectID, model.ItemID, getListItemByIDService.Options{
			ValueOnly: model.ResolvedOptions.ValueOnly,
		}), auth.NewAnonymousAuthentication())

		return request.SendPublicResponse[getListItemByIDService.Model](handler, c, http.StatusOK, nil, false)
	}
}
