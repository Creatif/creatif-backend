package project

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/app"
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/projects/createProject"
	"creatif/pkg/lib/logger"
	"github.com/labstack/echo/v4"
	"net/http"
)

func CreateProjectHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model app.CreateProject
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = app.SanitizeProject(model)

		l := logger.NewLogBuilder()
		handler := createProject.New(createProject.NewModel(model.Name), auth.NewFrontendAuthentication(), l)

		return request.SendResponse[createProject.Model](handler, c, http.StatusCreated, l, nil)
	}
}
