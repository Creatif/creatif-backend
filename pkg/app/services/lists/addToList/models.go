package addToList

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/locales"
	"creatif/pkg/app/services/shared/connections"
	"creatif/pkg/lib/constants"
	"creatif/pkg/lib/sdk"
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type VariableModel struct {
	Name      string   `json:"name"`
	Metadata  []byte   `json:"metadata"`
	Groups    []string `json:"groups"`
	Locale    string   `json:"locale"`
	Behaviour string   `json:"behaviour"`
	Value     []byte   `json:"value"`
}

type Model struct {
	Entry       VariableModel
	Name        string
	ProjectID   string
	Connections []connections.Connection
	ImagePaths  []string
}

type LogicModel struct {
	Variable    declarations.ListVariable
	Connections []declarations.Connection
	Groups      []string
}

func NewModel(projectId, name string, entry VariableModel, connections []connections.Connection, imagePaths []string) Model {
	return Model{
		Name:        name,
		ProjectID:   projectId,
		Entry:       entry,
		Connections: connections,
		ImagePaths:  imagePaths,
	}
}

func (a *Model) Validate() map[string]string {
	v := map[string]interface{}{
		"groups":           a.Entry.Groups,
		"name":             a.Name,
		"projectID":        a.ProjectID,
		"locale":           a.Entry.Locale,
		"behaviour":        a.Entry.Behaviour,
		"connectionsValid": nil,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("name", validation.Required),
			validation.Key("projectID", validation.Required, validation.RuneLength(27, 27)),
			validation.Key("locale", validation.Required, validation.By(func(value interface{}) error {
				t := value.(string)

				if !locales.ExistsByAlpha(t) {
					return errors.New(fmt.Sprintf("Locale '%s' does not exist.", t))
				}

				return nil
			})),
			validation.Key("behaviour", validation.Required, validation.By(func(value interface{}) error {
				v := value.(string)
				if v != constants.ReadonlyBehaviour && v != constants.ModifiableBehaviour {
					return errors.New(fmt.Sprintf("Invalid value for behaviour in variable '%s'. Variable behaviour can be 'modifiable' or 'readonly'", v))
				}

				return nil
			})),
			validation.Key("groups", validation.When(len(a.Entry.Groups) != 0, validation.Each(validation.RuneLength(1, 200))), validation.By(func(value interface{}) error {
				if a.Entry.Groups != nil {
					if len(a.Entry.Groups) > 20 {
						return errors.New("Maximum number of groups is 20.")
					}

					return nil
				}

				return nil
			})),
			validation.Key("connectionsValid", validation.By(func(value interface{}) error {
				if len(a.Connections) > 0 {
					names := make([]string, len(a.Connections))
					for _, r := range a.Connections {
						if r.StructureType != "map" && r.StructureType != "list" && r.StructureType != "variable" {
							return errors.New("Invalid connection structure type. Structure can can be one of: map, variable or list")
						}

						if r.Path == "" {
							return errors.New("Invalid connection. Path cannot be blank")
						}

						if sdk.Includes(names, r.Path) {
							return errors.New(fmt.Sprintf("Invalid connection. Duplicate connections are not possible. Structure type: %s; VariableID: %s", r.StructureType, r.VariableID))
						}

						names = append(names, r.Path)
					}

				}

				return nil
			})),
		),
	); err != nil {
		return sdk.ErrorToResponseError(err)
	}

	return nil
}
