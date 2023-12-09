package lists

import (
	declarations2 "creatif/cmd/http/handlers/declarations"
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations/lists"
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/lists/appendToList"
	"creatif/pkg/lib/logger"
	"creatif/pkg/lib/sdk"
	"github.com/labstack/echo/v4"
	"net/http"
)

func AppendToListHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model lists.AppendToList
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = lists.SanitizeAppendToList(model)
		if model.Locale == "" {
			model.Locale = declarations2.DefaultLocale
		}

		l := logger.NewLogBuilder()
		authentication := auth.NewApiAuthentication(request.GetApiAuthenticationCookie(c), l)
		handler := appendToList.New(appendToList.NewModel(
			model.ProjectID,
			model.Name,
			model.ID,
			model.ShortID,
			sdk.Map(model.Variables, func(idx int, value lists.AppendToListVariable) appendToList.Variable {
				return appendToList.Variable{
					Name:      value.Name,
					Metadata:  []byte(value.Metadata),
					Groups:    value.Groups,
					Locale:    "eng",
					Behaviour: value.Behaviour,
					Value:     []byte(value.Value),
				}
			}),
		), authentication, l)

		return request.SendResponse[appendToList.Model](handler, c, http.StatusCreated, l, func(c echo.Context, model interface{}) error {
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
