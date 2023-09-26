package declarations

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations"
	delete "creatif/pkg/app/declarations/removeMapEntry"
	"github.com/labstack/echo/v4"
	"net/http"
)

func DeleteMapEntry() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model declarations.DeleteMapEntry
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = declarations.SanitizeDeleteMapEntry(model)

		handler := delete.New(delete.NewModel(model.ProjectID, model.Name, model.EntryName))

		return request.SendResponse[delete.Model](handler, c, http.StatusCreated)
	}
}
