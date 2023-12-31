package addToMap

import (
	"creatif/pkg/app/services/locales"
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
	Behaviour string   `json:"behaviour"`
	Value     []byte   `json:"value"`
}

type Model struct {
	Entry     VariableModel
	Name      string
	ProjectID string
	Locale    string
}

func NewModel(projectId, locale, name string, entry VariableModel) Model {
	return Model{
		Name:      name,
		ProjectID: projectId,
		Locale:    locale,
		Entry:     entry,
	}
}

func (a *Model) Validate() map[string]string {
	v := map[string]interface{}{
		"groups":    a.Entry.Groups,
		"name":      a.Name,
		"projectID": a.ProjectID,
		"locale":    a.Locale,
		"behaviour": a.Entry.Behaviour,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("name", validation.When(a.Name != "", validation.RuneLength(1, 200))),
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
		),
	); err != nil {
		return sdk.ErrorToResponseError(err)
	}

	return nil
}
