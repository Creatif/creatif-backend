package lists

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations/lists"
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/lists/updateList"
	"github.com/labstack/echo/v4"
	"net/http"
)

func UpdateListHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model lists.UpdateList
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = lists.SanitizeUpdateList(model)

		authentication := auth.NewApiAuthentication(request.GetApiAuthenticationCookie(c))
		handler := updateList.New(updateList.NewModel(
			model.ProjectID,
			model.Fields,
			model.Name,
			model.Values.Name,
		), authentication)

		return request.SendResponse[updateList.Model](handler, c, http.StatusOK, func(c echo.Context, model interface{}) error {
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
