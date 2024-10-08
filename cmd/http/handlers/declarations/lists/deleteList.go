package lists

import (
	declarations2 "creatif/cmd/http/handlers/declarations"
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations/lists"
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/lists/deleteList"
	"github.com/labstack/echo/v4"
	"net/http"
)

func DeleteListHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model lists.DeleteList
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = lists.SanitizeDeleteList(model)
		if model.Locale == "" {
			model.Locale = declarations2.DefaultLocale
		}

		authentication := auth.NewApiAuthentication(request.GetApiAuthenticationCookie(c))
		handler := deleteList.New(deleteList.NewModel(
			model.ProjectID,
			model.Name,
			model.ID,
			model.ShortID,
		), authentication)

		return request.SendResponse[deleteList.Model](handler, c, http.StatusOK, func(c echo.Context, model interface{}) error {
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
