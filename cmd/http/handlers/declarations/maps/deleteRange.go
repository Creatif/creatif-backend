package maps

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations/maps"
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/maps/deleteRangeByID"
	"github.com/labstack/echo/v4"
	"net/http"
)

func DeleteRange() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model maps.DeleteRange
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = maps.SanitizeDeleteRange(model)

		a := auth.NewApiAuthentication(request.GetApiAuthenticationCookie(c))
		handler := deleteRangeByID.New(deleteRangeByID.NewModel(
			model.ProjectID,
			model.Name,
			model.Items,
		), a)

		return request.SendResponse[deleteRangeByID.Model](handler, c, http.StatusOK, func(c echo.Context, model interface{}) error {
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
