package structures

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/app"
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/structures/createAndDiff"
	"creatif/pkg/lib/sdk"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetStructureMetadataHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model app.GetStructureMetadata
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = app.SanitizeGetStructureMetadata(model)

		a := auth.NewApiAuthentication(request.GetApiAuthenticationCookie(c))
		handler := createAndDiff.New(createAndDiff.NewModel(model.ID, sdk.Map(model.Config, func(idx int, value app.GetStructureMetadataConfig) createAndDiff.Structure {
			return createAndDiff.Structure{
				Name: value.Name,
				Type: value.Type,
			}
		})), a)

		return request.SendResponse[createAndDiff.Model](handler, c, http.StatusOK, func(c echo.Context, model interface{}) error {
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
