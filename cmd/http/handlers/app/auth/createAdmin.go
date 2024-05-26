package auth

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/app"
	"creatif/pkg/app/services/auth/createAdmin"
	"creatif/pkg/lib/logger"
	"github.com/labstack/echo/v4"
	"net/http"
)

func CreateAdminHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model app.RegisterEmail
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = app.SanitizeRegisterEmail(model)

		l := logger.NewLogBuilder()
		handler := createAdmin.New(createAdmin.NewModel(model.Name, model.LastName, model.Email, model.Password), l)

		return request.SendResponse[createAdmin.Model](handler, c, http.StatusCreated, l, nil, false)
	}
}
