package queryMapVariable

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/sdk"
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Model struct {
	Name                    string
	ItemID                  string
	ProjectID               string
	ConnectionReplaceMethod string
}

func NewModel(projectId, name, itemID, connectionReplaceMethod string) Model {
	if connectionReplaceMethod == "" {
		connectionReplaceMethod = "connectionOnly"
	}

	return Model{
		ProjectID:               projectId,
		Name:                    name,
		ItemID:                  itemID,
		ConnectionReplaceMethod: connectionReplaceMethod,
	}
}

type LogicModel struct {
	Variable                  QueryVariable
	ChildConnectionStructures []ChildConnectionStructure
	Connections               []declarations.Connection
}

func (a *Model) Validate() map[string]string {
	v := map[string]interface{}{
		"name":                    a.Name,
		"itemId":                  a.ItemID,
		"projectID":               a.ProjectID,
		"connectionReplaceMethod": a.ConnectionReplaceMethod,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("name", validation.Required),
			validation.Key("itemId", validation.Required),
			validation.Key("projectID", validation.Required, validation.RuneLength(27, 27)),
			validation.Key("connectionReplaceMethod", validation.Required, validation.By(func(value interface{}) error {
				method := value.(string)

				if method == "" {
					return nil
				}

				if method != "fullReplacement" && method != "connectionOnly" {
					return errors.New("Connection replace method must be either empty or 'fullReplacement' and 'connectionOnly'. Default is 'connectionOnly'")
				}

				return nil
			})),
		),
	); err != nil {
		return sdk.ErrorToResponseError(err)
	}

	return nil
}
