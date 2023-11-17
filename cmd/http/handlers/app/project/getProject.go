package project

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/app"
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/projects/getProject"
	"creatif/pkg/lib/logger"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetProjectHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model app.GetProject
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = app.SanitizeGetProject(model)

		l := logger.NewLogBuilder()
		handler := getProject.New(getProject.NewModel(model.ProjectID), auth.NewFrontendAuthentication(request.GetAuthenticationCookie(c), l), l)

		return request.SendResponse[getProject.Model](handler, c, http.StatusOK, l, nil)
	}
}
