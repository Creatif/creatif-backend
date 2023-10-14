package switchByIndex

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
	Source      int64
	Destination int64
	Locale      string
	ProjectID   string
}

func NewModel(projectId, locale, name string, source, destination int64) Model {
	return Model{
		ProjectID:   projectId,
		Locale:      locale,
		Source:      source,
		Destination: destination,
		Name:        name,
	}
}

type LogicResult struct {
	To     declarations.ListVariable
	From   declarations.ListVariable
	Locale string
}

type ViewSourceDestination struct {
	ID      string `json:"id"`
	Index   string `json:"index"`
	ShortID string `json:"shortId"`
	Locale  string `json:"locale"`
	Name    string `json:"name"`
}

type View struct {
	Source      ViewSourceDestination `json:"newSource"`
	Destination ViewSourceDestination `json:"newDestination"`
}

func newView(model LogicResult) View {
	return View{
		Source: ViewSourceDestination{
			ID:      model.From.ID,
			Locale:  model.Locale,
			Index:   model.From.Index,
			ShortID: model.From.ShortID,
			Name:    model.From.Name,
		},
		Destination: ViewSourceDestination{
			ID:      model.To.ID,
			Index:   model.To.Index,
			Locale:  model.Locale,
			ShortID: model.To.ShortID,
			Name:    model.To.ShortID,
		},
	}
}

func (a *Model) Validate() map[string]string {
	v := map[string]interface{}{
		"name":      a.Name,
		"locale":    a.Locale,
		"projectID": a.ProjectID,
	}

	if err := validation.Validate(v,
		validation.Map(
			// Name cannot be empty, and the length must be between 5 and 20.
			validation.Key("name", validation.Required, validation.RuneLength(1, 200)),
			validation.Key("projectID", validation.Required, validation.RuneLength(26, 26)),
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
