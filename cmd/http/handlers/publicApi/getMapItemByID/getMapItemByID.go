package getMapItemByID

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/publicApi"
	"creatif/pkg/app/auth"
	getMapItemByIDService "creatif/pkg/app/services/publicApi/getMapItemById"
	"creatif/pkg/lib/logger"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetMapItemByIDHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model publicApi.GetMapItemByID
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = publicApi.SanitizeGetMapItemByID(model)

		l := logger.NewLogBuilder()
		handler := getMapItemByIDService.New(getMapItemByIDService.NewModel(model.ProjectID, model.ItemID), auth.NewAnonymousAuthentication(), l)

		return request.SendPublicResponse[getMapItemByIDService.Model](handler, c, http.StatusOK, l, nil, false)
	}
}
