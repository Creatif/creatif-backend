package updateMapVariable

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/constants"
	"creatif/pkg/lib/sdk"
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"time"
)

type VariableModel struct {
	Name      string   `json:"name"`
	Metadata  []byte   `json:"metadata"`
	Groups    []string `json:"groups"`
	Behaviour string   `json:"behaviour"`
	Value     []byte   `json:"value"`
}

type Model struct {
	Entry     VariableModel `json:"entry"`
	Name      string        `json:"name"`
	ProjectID string        `query:"projectID"`
}

func NewModel(projectId, name string, entry VariableModel) Model {
	return Model{
		Name:      name,
		ProjectID: projectId,
		Entry:     entry,
	}
}

type LogicResult struct {
	Map   declarations.Map
	Entry declarations.MapVariable
}

func (a *Model) Validate() map[string]string {
	v := map[string]interface{}{
		"groups":    a.Entry.Groups,
		"name":      a.Name,
		"behaviour": a.Entry.Behaviour,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("name", validation.Required, validation.RuneLength(1, 200)),
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

type ViewEntry struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	Metadata  interface{} `json:"metadata"`
	Groups    []string    `json:"groups"`
	Behaviour string      `json:"behaviour"`
	Value     interface{} `json:"value"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
}

type View struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	ProjectID string `json:"projectID"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Entry ViewEntry `json:"entry"`
}

func newView(logicResult LogicResult) View {
	m := logicResult.Map
	variable := logicResult.Entry

	return View{
		ID:        m.ID,
		Name:      m.Name,
		ProjectID: m.ProjectID,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		Entry: ViewEntry{
			ID:        variable.ID,
			Name:      variable.Name,
			Metadata:  variable.Metadata,
			Groups:    variable.Groups,
			Behaviour: variable.Behaviour,
			Value:     variable.Value,
			CreatedAt: variable.CreatedAt,
			UpdatedAt: variable.UpdatedAt,
		},
	}
}
