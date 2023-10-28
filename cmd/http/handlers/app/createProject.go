package app

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/app"
	create "creatif/pkg/app/app/createProject"
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
		handler := create.New(create.NewModel(model.Name), l)

		return request.SendResponse[create.Model](handler, c, http.StatusCreated, l)
	}
}
