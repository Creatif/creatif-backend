package deleteList

import (
	"creatif/pkg/lib/sdk"
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Model struct {
	Name      string
	ShortID   string
	ID        string
	ProjectID string
}

func NewModel(projectId, name, id, shortID string) Model {
	return Model{
		Name:      name,
		ID:        id,
		ShortID:   shortID,
		ProjectID: projectId,
	}
}

func (a Model) Validate() map[string]string {
	v := map[string]interface{}{
		"name":      a.Name,
		"id":        a.ID,
		"idExists":  nil,
		"projectID": a.ProjectID,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("name", validation.When(a.Name != "", validation.RuneLength(1, 200))),
			validation.Key("id", validation.When(a.ID != "", validation.RuneLength(27, 27))),
			validation.Key("idExists", validation.By(func(value interface{}) error {
				name := a.Name
				shortId := a.ShortID
				id := a.ID

				if name == "" && shortId == "" && id == "" {
					return errors.New("At least one of 'id', 'name' or 'shortID' must be supplied in order to identify this variable.")
				}
				return nil
			})), validation.Key("projectID", validation.Required, validation.RuneLength(1, 27)),
		),
	); err != nil {
		return sdk.ErrorToResponseError(err)
	}

	return nil
}
