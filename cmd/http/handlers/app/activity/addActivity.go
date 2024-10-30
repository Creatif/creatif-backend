package activity

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/app"
	"creatif/pkg/app/auth"
	createActivity "creatif/pkg/app/services/activity/create"
	"github.com/labstack/echo/v4"
	"net/http"
)

func AddActivityHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model app.AddActivity
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = app.SanitizeAddActivity(model)
		authentication := auth.NewApiAuthentication(request.GetApiAuthenticationCookie(c))
		handler := createActivity.New(createActivity.NewModel(model.ProjectID, []byte(model.Data)), authentication)

		return request.SendResponse[createActivity.Model](handler, c, http.StatusCreated, func(c echo.Context, model interface{}) error {
			if authentication.ShouldRefresh() {
				session, err := authentication.Refresh()
				if err != nil {
					return err
				}

				c.SetCookie(request.EncryptAuthenticationCookie(session))
			}

			return nil
		}, false)
	}
}
