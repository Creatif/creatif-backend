package appendToList

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
		"variableLen":     len(a.Variables),
		"projectID":       a.ProjectID,
		"name":            a.Name,
		"variableLocales": nil,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("name", validation.Required),
			validation.Key("variableLocales", validation.By(func(value interface{}) error {
				for _, v := range a.Variables {
					l := v.Locale
					if l != "" && !locales.ExistsByAlpha(l) {
						return errors.New(fmt.Sprintf("Locale '%s' does not exist for variable with name '%s'", l, v.Name))
					}
				}

				return nil
			})),
			validation.Key("projectID", validation.Required, validation.RuneLength(26, 26)),
			validation.Key("variableLen", validation.By(func(value interface{}) error {
				l := value.(int)

				if l > 1000 {
					return errors.New("The number of variables when creating a list cannot be higher than 1000.")
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
	ProjectID string `json:"projectID"`
	Name      string `json:"name"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func newView(model declarations.List) View {
	return View{
		ID:        model.ID,
		ProjectID: model.ProjectID,
		Name:      model.Name,

		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}
