package mapCreate

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/constants"
	"creatif/pkg/lib/sdk"
	"creatif/pkg/lib/storage"
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gorm.io/gorm"
)

type VariableModel struct {
	Name      string   `json:"name"`
	Metadata  []byte   `json:"metadata"`
	Groups    []string `json:"groups"`
	Behaviour string   `json:"behaviour"`
	Value     []byte   `json:"value"`
}

type View struct {
	ID        string              `json:"id"`
	Name      string              `json:"name"`
	Variables []map[string]string `json:"variables"`
}

type Entry struct {
	Type  string
	Model interface{}
}

type Model struct {
	Entries   []Entry `json:"entries"`
	Name      string  `json:"name"`
	ProjectID string  `json:"projectID"`
}

type LogicResult struct {
	ID        string
	Variables []map[string]string
	Name      string
}

func NewModel(projectId, name string, entries []Entry) Model {
	return Model{
		Name:      name,
		ProjectID: projectId,
		Entries:   entries,
	}
}

func (a *Model) Validate() map[string]string {
	v := map[string]interface{}{
		"name":               a.Name,
		"uniqueName":         a.Name,
		"validNum":           a.Entries,
		"validVariableNames": a.Entries,
		"behaviourValid":     a.Entries,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("name", validation.Required, validation.RuneLength(1, 200)),
			validation.Key("uniqueName", validation.By(func(value interface{}) error {
				name := value.(string)

				var m declarations.Map
				if err := storage.GetBy((&declarations.Map{}).TableName(), "name", name, &m); err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
					return errors.New("Map with this name already exists")
				}

				return nil
			})),
			validation.Key("validNum", validation.By(func(value interface{}) error {
				if len(a.Entries) == 0 {
					return errors.New("Empty entries are not permitted. Maps must have values.")
				}

				if len(a.Entries) > 1000 {
					return errors.New("Number of map values cannot be larger than 1000.")
				}

				return nil
			})),
			validation.Key("validVariableNames", validation.By(func(value interface{}) error {
				m := make(map[string]int)
				for _, entry := range a.Entries {
					if entry.Type == "variable" {
						o := entry.Model.(VariableModel)
						m[o.Name] = 0
					}
				}

				if len(m) != len(a.Entries) {
					return errors.New("Some variable/map names are not unique. All variable/map names must be unique.")
				}
				return nil
			})),
			validation.Key("behaviourValid", validation.Required, validation.By(func(value interface{}) error {
				m := make(map[string]string)
				for _, entry := range a.Entries {
					if entry.Type == "variable" {
						o := entry.Model.(VariableModel)
						m[o.Name] = o.Behaviour
					}
				}

				for key, v := range m {
					if v != constants.ReadonlyBehaviour && v != constants.ModifiableBehaviour {
						return errors.New(fmt.Sprintf("Invalid value for behaviour in variable '%s'. Variable behaviour can be 'modifiable' or 'readonly'", key))
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

func newView(model LogicResult) View {
	return View{
		ID:        model.ID,
		Name:      model.Name,
		Variables: model.Variables,
	}
}
