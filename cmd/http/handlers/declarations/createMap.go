package declarations

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations"
	"creatif/pkg/app/declarations/mapCreate"
	"github.com/labstack/echo/v4"
	"net/http"
)

func CreateMapHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model declarations.CreateMap
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = declarations.SanitizeMapModel(model)

		serviceEntries := make([]mapCreate.Entry, 0)
		for _, entry := range model.Entries {

			m, ok := entry.Model.(declarations.MapVariableModel)
			if ok {
				serviceEntries = append(serviceEntries, mapCreate.Entry{
					Type: entry.Type,
					Model: mapCreate.VariableModel{
						Name:      m.Name,
						Metadata:  m.Metadata,
						Groups:    m.Groups,
						Behaviour: m.Behaviour,
					},
				})
			}
		}

		handler := mapCreate.New(mapCreate.NewModel(model.Name, serviceEntries))

		return request.SendResponse[mapCreate.Model](handler, c, http.StatusCreated)
	}
}
