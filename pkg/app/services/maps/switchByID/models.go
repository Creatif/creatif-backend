package switchByID

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/locales"
	"creatif/pkg/lib/sdk"
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"strings"
)

type Model struct {
	// this can be project name
	Name           string
	Source         string
	Destination    string
	ProjectID      string
	OrderDirection string
}

func NewModel(projectId, name, source, destination, orderDirection string) Model {
	return Model{
		ProjectID:      projectId,
		Source:         source,
		Destination:    destination,
		Name:           name,
		OrderDirection: orderDirection,
	}
}

type LogicResult struct {
	To   declarations.ListVariable
	From declarations.ListVariable
}

type ViewSourceDestination struct {
	ID        string   `json:"id"`
	Locale    string   `json:"locale"`
	ShortID   string   `json:"shortId"`
	Name      string   `json:"name"`
	Behaviour string   `json:"behaviour"`
	Groups    []string `json:"groups"`
}

type View struct {
	Source      ViewSourceDestination `json:"source"`
	Destination ViewSourceDestination `json:"destination"`
}

func newView(model LogicResult) View {
	sourceLocale, _ := locales.GetAlphaWithID(model.To.LocaleID)
	destinationLocale, _ := locales.GetAlphaWithID(model.From.LocaleID)
	return View{
		Source: ViewSourceDestination{
			ID:        model.From.ID,
			ShortID:   model.From.ShortID,
			Name:      model.From.Name,
			Locale:    destinationLocale,
			Behaviour: model.From.Behaviour,
		},
		Destination: ViewSourceDestination{
			ID:        model.To.ID,
			ShortID:   model.To.ShortID,
			Locale:    sourceLocale,
			Name:      model.To.Name,
			Behaviour: model.To.Behaviour,
		},
	}
}

func (a *Model) Validate() map[string]string {
	v := map[string]interface{}{
		"name":           a.Name,
		"projectID":      a.ProjectID,
		"source":         a.Source,
		"destination":    a.Destination,
		"orderDirection": a.OrderDirection,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("name", validation.Required),
			validation.Key("projectID", validation.Required, validation.RuneLength(27, 27)),
			validation.Key("source", validation.Required, validation.RuneLength(27, 27)),
			validation.Key("destination", validation.Required, validation.RuneLength(27, 27)),
			validation.Key("orderDirection", validation.Required, validation.By(func(value interface{}) error {
				direction := strings.ToLower(value.(string))

				if direction != "asc" && direction != "desc" {
					return errors.New("Invalid order direction. Order direction can be either 'DESC' or 'ASC'")
				}
				return nil
			})),
		),
	); err != nil {
		return sdk.ErrorToResponseError(err)
	}

	return nil
}
