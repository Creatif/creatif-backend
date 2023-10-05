package declarations

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations"
	mapCreate2 "creatif/pkg/app/services/maps/mapCreate"
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

		serviceEntries := make([]mapCreate2.Entry, 0)
		for _, entry := range model.Entries {
			m, ok := entry.Model.(declarations.MapVariableModel)
			if ok {
				serviceEntries = append(serviceEntries, mapCreate2.Entry{
					Type: entry.Type,
					Model: mapCreate2.VariableModel{
						Name:      m.Name,
						Metadata:  []byte(m.Metadata),
						Value:     []byte(m.Value),
						Groups:    m.Groups,
						Behaviour: m.Behaviour,
					},
				})
			}
		}

		handler := mapCreate2.New(mapCreate2.NewModel(model.ProjectID, model.Name, serviceEntries))

		return request.SendResponse[mapCreate2.Model](handler, c, http.StatusCreated)
	}
}
