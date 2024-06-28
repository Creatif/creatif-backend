package images

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations/images"
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/images/getImage"
	"creatif/pkg/lib/logger"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetImageHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model images.GetImage
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = images.SanitizeGetImage(model)

		l := logger.NewLogBuilder()
		authentication := auth.NewApiAuthentication(request.GetApiAuthenticationCookie(c), l)
		handler := getImage.New(getImage.NewModel(model.ProjectID, model.StructureID), authentication, l)

		img, err := handler.Handle()
		if err != nil {
			return c.NoContent(http.StatusNotFound)
		}

		return c.File(img.Name)
	}
}
