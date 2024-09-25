package lists

import (
	declarations2 "creatif/cmd/http/handlers/declarations"
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations/lists"
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/lists/deleteListItemByID"
	"github.com/labstack/echo/v4"
	"net/http"
)

func DeleteListItemByIDHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model lists.DeleteListItemByID
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = lists.SanitizeDeleteListItemByID(model)
		if model.Locale == "" {
			model.Locale = declarations2.DefaultLocale
		}

		a := auth.NewApiAuthentication(request.GetApiAuthenticationCookie(c))
		handler := deleteListItemByID.New(deleteListItemByID.NewModel(
			model.ProjectID,
			model.Name,
			model.ItemID,
		), a)

		return request.SendResponse[deleteListItemByID.Model](handler, c, http.StatusOK, func(c echo.Context, model interface{}) error {
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
