package queryMapVariable

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/sdk"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Model struct {
	Name      string
	ItemID    string
	ProjectID string
}

func NewModel(projectId, name, itemID string) Model {
	return Model{
		ProjectID: projectId,
		Name:      name,
		ItemID:    itemID,
	}
}

type LogicModel struct {
	Variable                  QueryVariable
	ChildConnectionStructures []ChildConnectionStructure
	Connections               []declarations.Connection
}

func (a *Model) Validate() map[string]string {
	v := map[string]interface{}{
		"name":      a.Name,
		"itemId":    a.ItemID,
		"projectID": a.ProjectID,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("name", validation.Required),
			validation.Key("itemId", validation.Required),
			validation.Key("projectID", validation.Required, validation.RuneLength(27, 27)),
		),
	); err != nil {
		return sdk.ErrorToResponseError(err)
	}

	return nil
}
