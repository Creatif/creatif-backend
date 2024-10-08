package lists

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations/lists"
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/lists/deleteRangeByID"
	"github.com/labstack/echo/v4"
	"net/http"
)

func DeleteRangeByIDHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model lists.DeleteRangeByID
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = lists.SanitizeDeleteRangeByID(model)

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
