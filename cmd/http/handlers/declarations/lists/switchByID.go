package lists

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations/lists"
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/lists/switchByID"
	"github.com/labstack/echo/v4"
	"net/http"
)

func SwitchByIDHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model lists.SwitchByID
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = lists.SanitizeSwitchByID(model)

		authentication := auth.NewApiAuthentication(request.GetApiAuthenticationCookie(c))
		handler := switchByID.New(switchByID.NewModel(
			model.ProjectID,
			model.Name,
			model.Source,
			model.Destination,
		), authentication)

		return request.SendResponse[switchByID.Model](handler, c, http.StatusOK, func(c echo.Context, model interface{}) error {
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
