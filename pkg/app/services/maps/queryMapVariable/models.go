package queryMapVariable

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/sdk"
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Model struct {
	Name               string
	ItemID             string
	ProjectID          string
	ConnectionViewType string
}

func NewModel(projectId, name, itemID, connectionViewType string) Model {
	if connectionViewType == "" {
		connectionViewType = "connection"
	}

	return Model{
		ProjectID:          projectId,
		Name:               name,
		ItemID:             itemID,
		ConnectionViewType: connectionViewType,
	}
}

type LogicModel struct {
	Variable                  QueryVariable
	ChildConnectionStructures []ChildConnectionStructure
	Connections               []declarations.Connection
}

func (a *Model) Validate() map[string]string {
	v := map[string]interface{}{
		"name":               a.Name,
		"itemId":             a.ItemID,
		"projectID":          a.ProjectID,
		"connectionViewType": a.ConnectionViewType,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("name", validation.Required),
			validation.Key("itemId", validation.Required),
			validation.Key("projectID", validation.Required, validation.RuneLength(27, 27)),
			validation.Key("connectionViewType", validation.Required, validation.By(func(value interface{}) error {
				method := value.(string)

				if method == "" {
					return nil
				}

				if method != "connection" && method != "value" && method != "variable" {
					return errors.New("Connection replace method must be either empty or 'connection', 'value' or 'variable'. Default is 'connection'")
				}

				return nil
			})),
		),
	); err != nil {
		return sdk.ErrorToResponseError(err)
	}

	return nil
}
