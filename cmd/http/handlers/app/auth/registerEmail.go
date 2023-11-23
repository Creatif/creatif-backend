package auth

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/app"
	"creatif/pkg/app/services/auth/registerEmail"
	"creatif/pkg/lib/logger"
	"github.com/labstack/echo/v4"
	"net/http"
)

func CreateRegisterEmailHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model app.RegisterEmail
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = app.SanitizeRegisterEmail(model)

		l := logger.NewLogBuilder()
		handler := registerEmail.New(registerEmail.NewModel(model.Name, model.LastName, model.Email, model.Password, model.PolicyAccepted), nil, l)

		return request.SendResponse[registerEmail.Model](handler, c, http.StatusCreated, l, nil, false)
	}
}
