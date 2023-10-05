package queryListByID

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/sdk"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/lib/pq"
	"time"
)

type Model struct {
	// this can be project name
	Name      string `json:"name"`
	ID        string `json:"id"`
	ProjectID string `json:"projectID"`
}

func NewModel(projectId, name, id string) Model {
	return Model{
		ProjectID: projectId,
		ID:        id,
		Name:      name,
	}
}

type View struct {
	ID        string         `json:"id"`
	Index     int64          `json:"index"`
	ShortID   string         `json:"shortId"`
	Name      string         `json:"name"`
	Behaviour string         `json:"behaviour"`
	Groups    pq.StringArray `json:"groups"`
	Metadata  interface{}    `json:"metadata"`
	Value     interface{}    `json:"value"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func newView(model declarations.ListVariable) View {
	return View{
		ID:        model.ID,
		Index:     model.Index,
		ShortID:   model.ShortID,
		Name:      model.Name,
		Behaviour: model.Behaviour,
		Groups:    model.Groups,
		Metadata:  model.Metadata,
		Value:     model.Value,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
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
