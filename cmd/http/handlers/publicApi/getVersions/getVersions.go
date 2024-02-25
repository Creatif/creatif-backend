package getVersions

import (
	"creatif/cmd/http/request"
	getVersionsRequest "creatif/cmd/http/request/publicApi/getVersions"
	"creatif/pkg/app/auth"
	publicApiGetVersions "creatif/pkg/app/services/publicApi/getVersions"
	"creatif/pkg/lib/logger"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetVersionsHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model getVersionsRequest.GetVersions
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = getVersionsRequest.SanitizeGetVersions(model)

		l := logger.NewLogBuilder()
		handler := publicApiGetVersions.New(publicApiGetVersions.NewModel(model.ProjectID), auth.NewAnonymousAuthentication(), l)

		return request.SendResponse[publicApiGetVersions.Model](handler, c, http.StatusCreated, l, nil, false)
	}
}
