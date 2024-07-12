package getFile

import (
	"creatif/cmd/http/request/publicApi"
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/publicApi/getFile"
	"creatif/pkg/lib/constants"
	"creatif/pkg/lib/logger"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"path/filepath"
)

func GetFileHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model publicApi.GetFile
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = publicApi.SanitizeGetFile(model)

		l := logger.NewLogBuilder()
		handler := getFile.New(getFile.NewModel(model.ProjectID, model.FileID, model.Version), auth.NewAnonymousAuthentication(), l)

		file, err := handler.Handle()
		if err != nil {
			return c.NoContent(http.StatusNotFound)
		}

		c.Set("Content-Type", file.MimeType)
		return c.File(fmt.Sprintf("%s/%s/%s/%s", constants.PublicDirectory, model.ProjectID, model.Version, filepath.Base(file.Name)))
	}
}
