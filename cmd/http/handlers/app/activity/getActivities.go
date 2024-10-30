package activity

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/app"
	"creatif/pkg/app/auth"
	getActivities "creatif/pkg/app/services/activity/get"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetActivityHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model app.GetActivities
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = app.SanitizeGetActivities(model)

		handler := getActivities.New(getActivities.NewModel(model.ProjectID), auth.NewApiAuthentication(request.GetAuthenticationCookie(c)))

		return request.SendResponse[getActivities.Model](handler, c, http.StatusOK, nil, false)
	}
}
