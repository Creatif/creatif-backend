package addToMap

import (
	"creatif/pkg/app/services/locales"
	"creatif/pkg/app/services/shared"
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
	Entry      VariableModel
	Name       string
	ProjectID  string
	References []shared.Reference
}

func NewModel(projectId, name string, entry VariableModel, references []shared.Reference) Model {
	return Model{
		Name:       name,
		ProjectID:  projectId,
		Entry:      entry,
		References: references,
	}
}

func (a *Model) Validate() map[string]string {
	v := map[string]interface{}{
		"groups":          a.Entry.Groups,
		"name":            a.Name,
		"projectID":       a.ProjectID,
		"locale":          a.Entry.Locale,
		"behaviour":       a.Entry.Behaviour,
		"referencesValid": nil,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("name", validation.Required),
			validation.Key("projectID", validation.Required, validation.RuneLength(26, 26)),
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
			validation.Key("referencesValid", validation.By(func(value interface{}) error {
				if len(a.References) > 0 {
					names := make([]string, len(a.References))
					for _, r := range a.References {
						if r.StructureType != "map" && r.StructureType != "list" && r.StructureType != "variable" {
							return errors.New("Invalid reference structure type. Structure can can be one of: map, variable or list")
						}

						name := fmt.Sprintf("%s_%s_%s", r.StructureName, r.StructureType, r.VariableID)
						if sdk.Includes(names, name) {
							return errors.New(fmt.Sprintf("Invalid reference. Duplicate reference are not possible. Structure name: %s; Structure type: %s; VariableID: %s", r.StructureName, r.StructureType, r.VariableID))
						}

						names = append(names, name)
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
