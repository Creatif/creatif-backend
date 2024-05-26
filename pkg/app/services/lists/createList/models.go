package createList

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/locales"
	"creatif/pkg/lib/sdk"
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"time"
)

type Variable struct {
	Name      string
	Metadata  []byte
	Groups    []string
	Behaviour string
	Locale    string
	Value     []byte
}

type Model struct {
	Name      string
	ProjectID string
	Variables []Variable
}

func NewModel(projectId, name string, variables []Variable) Model {
	return Model{
		Name:      name,
		ProjectID: projectId,
		Variables: variables,
	}
}

func (a Model) Validate() map[string]string {
	v := map[string]interface{}{
		"name":            a.Name,
		"projectID":       a.ProjectID,
		"variableLen":     len(a.Variables),
		"variableLocales": nil,
		"groups":          nil,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("projectID", validation.Required, validation.RuneLength(27, 27)),
			validation.Key("name", validation.Required, validation.RuneLength(1, 200)),
			validation.Key("variableLen", validation.By(func(value interface{}) error {
				l := value.(int)

				if l > 1000 {
					return errors.New("The number of variables when creating a list cannot be higher than 1000.")
				}

				return nil
			})),
			validation.Key("variableLocales", validation.By(func(value interface{}) error {
				for _, v := range a.Variables {
					l := v.Locale
					if !locales.ExistsByAlpha(l) {
						return errors.New(fmt.Sprintf("Locale '%s' does not exist for variable with name '%s'", l, v.Name))
					}
				}

				return nil
			})),
			validation.Key("groups", validation.By(func(value interface{}) error {
				for _, variable := range a.Variables {
					if len(variable.Groups) > 20 {
						return errors.New(fmt.Sprintf("Invalid number of groups for '%s'. Maximum number of groups per variable is 20.", variable.Name))
					}

					for _, g := range variable.Groups {
						if len(g) > 100 {
							return errors.New(fmt.Sprintf("Invalid group length for '%s'. Maximum number of characters per groups is 200.", g))
						}
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

type View struct {
	ID        string `json:"id"`
	ShortID   string `json:"shortID"`
	ProjectID string `json:"projectID"`
	Name      string `json:"name"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func newView(model declarations.List) View {
	return View{
		ID:        model.ID,
		ShortID:   model.ShortID,
		ProjectID: model.ProjectID,
		Name:      model.Name,

		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}
