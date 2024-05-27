package truncateStructure

import (
	"creatif/pkg/lib/sdk"
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Model struct {
	ID        string `json:"id"`
	ProjectID string `json:"projectId"`
	Type      string `json:"type"`
}

func NewModel(projectId, id, t string) Model {
	return Model{
		ID:        id,
		ProjectID: projectId,
		Type:      t,
	}
}

func (a Model) Validate() map[string]string {
	v := map[string]interface{}{
		"projectId": a.ProjectID,
		"id":        a.ID,
		"type":      a.Type,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("projectId", validation.Required, validation.RuneLength(27, 27)),
			validation.Key("id", validation.Required, validation.RuneLength(27, 27)),
			validation.Key("type", validation.By(func(value interface{}) error {
				t := value.(string)

				if t != "map" && t != "list" {
					return errors.New("Type must be either 'map' or 'list'")
				}

				return nil
			})),
		),
	); err != nil {
		return sdk.ErrorToResponseError(err)
	}

	return nil
}
