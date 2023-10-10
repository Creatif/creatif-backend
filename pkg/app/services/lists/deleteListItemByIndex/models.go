package deleteListItemByIndex

import (
	"creatif/pkg/lib/sdk"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Model struct {
	Name      string
	ItemIndex string
	ProjectID string
}

func NewModel(projectId, name, itemIndex string) Model {
	return Model{
		Name:      name,
		ItemIndex: itemIndex,
		ProjectID: projectId,
	}
}

func (a Model) Validate() map[string]string {
	v := map[string]interface{}{
		"name":      a.Name,
		"projectID": a.ProjectID,
		"itemIndex": a.ItemIndex,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("name", validation.Required, validation.RuneLength(1, 200)),
			validation.Key("projectID", validation.Required, validation.RuneLength(1, 26)),
			validation.Key("itemIndex", validation.Required, validation.RuneLength(1, 26)),
		),
	); err != nil {
		return sdk.ErrorToResponseError(err)
	}

	return nil
}
