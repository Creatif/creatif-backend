package create

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/sdk"
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"time"
)

type GetNodeModel struct {
	// this can be project name
	ID string `json:"id"`
	// TODO: Add project ID prop here
}

func NewGetNodeModel(id string) GetNodeModel {
	return GetNodeModel{
		ID: id,
	}
}

type View struct {
	ID        string                 `json:"id"`
	Name      string                 `json:"name"`
	Type      string                 `json:"type"`
	Groups    []string               `json:"groups"`
	Behaviour string                 `json:"behaviour"`
	Metadata  map[string]interface{} `json:"metadata"`

	CreatedAt time.Time `gorm:"<-:create" json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func newView(model declarations.Node) View {
	return View{
		ID:        model.ID,
		Name:      model.Name,
		Type:      model.Type,
		Groups:    model.Groups,
		Behaviour: model.Behaviour,
		Metadata:  sdk.UnmarshalToMap([]byte(model.Metadata)),
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}

func (a *GetNodeModel) Validate() map[string]string {
	v := map[string]interface{}{
		"id": a.ID,
	}

	if err := validation.Validate(v,
		validation.Map(
			// Name cannot be empty, and the length must be between 5 and 20.
			validation.Key("id", validation.Required),
		),
	); err != nil {
		var e map[string]string
		b, err := json.Marshal(err)
		if err != nil {
			return map[string]string{
				"unrecoverable": "An internal validation error occurred. This should not happen. Please, submit a bug.",
			}
		}

		if err := json.Unmarshal(b, &e); err != nil {
			return map[string]string{
				"unrecoverable": "An internal validation error occurred. This should not happen. Please, submit a bug.",
			}
		}

		return e
	}

	return nil
}
