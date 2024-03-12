package getListItemByID

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/publicApi"
	"creatif/pkg/app/auth"
	getListItemByIDService "creatif/pkg/app/services/publicApi/getListItemById"
	"creatif/pkg/lib/logger"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetListItemByIDHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model publicApi.GetListItemByID
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = publicApi.SanitizeGetListItemByID(model)

		l := logger.NewLogBuilder()
		handler := getListItemByIDService.New(getListItemByIDService.NewModel(model.ProjectID, model.ItemID), auth.NewAnonymousAuthentication(), l)

		return request.SendPublicResponse[getListItemByIDService.Model](handler, c, http.StatusOK, l, nil, false)
	}
}
