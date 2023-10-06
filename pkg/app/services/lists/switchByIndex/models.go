package switchByIndex

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/sdk"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Model struct {
	// this can be project name
	Name        string `json:"name"`
	Source      int64  `json:"fromIndex"`
	Destination int64  `json:"toIndex"`
	ProjectID   string `json:"projectID"`
}

func NewModel(projectId, name string, source, destination int64) Model {
	return Model{
		ProjectID:   projectId,
		Source:      source,
		Destination: destination,
		Name:        name,
	}
}

type LogicResult struct {
	To   declarations.ListVariable
	From declarations.ListVariable
}

type ViewSourceDestination struct {
	ID      string `json:"id"`
	Index   string `json:"index"`
	ShortID string `json:"shortId"`
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
			Index:   model.From.Index,
			ShortID: model.From.ShortID,
			Name:    model.From.Name,
		},
		Destination: ViewSourceDestination{
			ID:      model.To.ID,
			Index:   model.To.Index,
			ShortID: model.To.ShortID,
			Name:    model.To.ShortID,
		},
	}
}

func (a *Model) Validate() map[string]string {
	v := map[string]interface{}{
		"name": a.Name,
	}

	if err := validation.Validate(v,
		validation.Map(
			// Name cannot be empty, and the length must be between 5 and 20.
			validation.Key("name", validation.Required, validation.RuneLength(1, 200)),
		),
	); err != nil {
		return sdk.ErrorToResponseError(err)
	}

	return nil
}
