package switchByID

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/locales"
	"creatif/pkg/lib/sdk"
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Model struct {
	// this can be project name
	Name        string
	ID          string
	ShortID     string
	Source      string
	Destination string
	ProjectID   string
	Locale      string
}

func NewModel(projectId, locale, name, id, shortID, source, destination string) Model {
	return Model{
		ProjectID:   projectId,
		Locale:      locale,
		Source:      source,
		Destination: destination,
		Name:        name,
		ID:          id,
		ShortID:     shortID,
	}
}

type LogicResult struct {
	To     declarations.ListVariable
	From   declarations.ListVariable
	Locale string
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
	return View{
		Source: ViewSourceDestination{
			ID:        model.From.ID,
			Locale:    model.Locale,
			ShortID:   model.From.ShortID,
			Name:      model.From.Name,
			Behaviour: model.From.Behaviour,
			Groups:    model.From.Groups,
		},
		Destination: ViewSourceDestination{
			ID:        model.To.ID,
			Locale:    model.Locale,
			ShortID:   model.To.ShortID,
			Name:      model.To.Name,
			Behaviour: model.To.Behaviour,
			Groups:    model.To.Groups,
		},
	}
}

func (a *Model) Validate() map[string]string {
	v := map[string]interface{}{
		"name":        a.Name,
		"id":          a.ID,
		"idExists":    nil,
		"projectID":   a.ProjectID,
		"source":      a.Source,
		"destination": a.Destination,
		"locale":      a.Locale,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("name", validation.When(a.Name != "", validation.RuneLength(1, 200))),
			validation.Key("id", validation.When(a.ID != "", validation.RuneLength(26, 26))),
			validation.Key("idExists", validation.By(func(value interface{}) error {
				name := a.Name
				shortId := a.ShortID
				id := a.ID

				if id != "" && len(id) != 26 {
					return errors.New("ID must have 26 characters")
				}

				if name == "" && shortId == "" && id == "" {
					return errors.New("At least one of 'id', 'name' or 'shortID' must be supplied in order to identify this list.")
				}
				return nil
			})),
			validation.Key("projectID", validation.Required, validation.RuneLength(26, 26)),
			validation.Key("source", validation.Required, validation.RuneLength(26, 26)),
			validation.Key("destination", validation.Required, validation.RuneLength(26, 26)),
			validation.Key("locale", validation.Required, validation.By(func(value interface{}) error {
				t := value.(string)

				if !locales.ExistsByAlpha(t) {
					return errors.New(fmt.Sprintf("Locale '%s' not found.", t))
				}

				return nil
			})),
		),
	); err != nil {
		return sdk.ErrorToResponseError(err)
	}

	return nil
}
