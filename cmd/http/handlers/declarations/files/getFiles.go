package files

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations/files"
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/files/getFile"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetFileHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model files.GetFile
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = files.SanitizeGetFile(model)

		authentication := auth.NewApiAuthentication(request.GetApiAuthenticationCookie(c))
		handler := getFile.New(getFile.NewModel(model.ProjectID, model.StructureID), authentication)

		img, err := handler.Handle()
		if err != nil {
			return c.NoContent(http.StatusNotFound)
		}

		return c.File(img.Name)
	}
}
