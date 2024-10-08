package toggleProduction

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/publishing/toggleProduction"
	"creatif/pkg/app/auth"
	toggleProductionService "creatif/pkg/app/services/publishing/toggleProduction"
	"github.com/labstack/echo/v4"
	"net/http"
)

func ToggleProductionHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model toggleProduction.ToggleProduction
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = toggleProduction.SanitizeToggleProduction(model)

		authentication := auth.NewApiAuthentication(request.GetApiAuthenticationCookie(c))
		handler := toggleProductionService.New(toggleProductionService.NewModel(model.ProjectID, model.ID), authentication)

		return request.SendResponse[toggleProductionService.Model](handler, c, http.StatusCreated, func(c echo.Context, model interface{}) error {
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
