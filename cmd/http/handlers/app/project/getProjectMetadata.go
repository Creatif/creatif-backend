package project

import (
	"creatif/cmd"
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/app"
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/projects/getProjectMetadata"
	"creatif/pkg/lib/logger"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetProjectMetadataHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model app.GetProject
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = app.SanitizeGetProject(model)

		apiKey := c.Request().Header.Get(cmd.CreatifApiHeader)
		projectId := c.Request().Header.Get(cmd.CreatifProjectIDHeader)

		l := logger.NewLogBuilder()
		a := auth.NewApiAuthentication(request.GetApiAuthenticationCookie(c), projectId, apiKey, l)
		handler := getProjectMetadata.New(a, l)

		return request.SendResponse(handler, c, http.StatusOK, l, func(c echo.Context, model interface{}) error {
			if a.ShouldRefresh() {
				session, err := a.Refresh()
				if err != nil {
					return err
				}

				c.SetCookie(request.EncryptAuthenticationCookie(session))
			}

			return nil
		}, false)
	}
}
