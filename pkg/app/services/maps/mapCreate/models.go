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
	Locale    string   `json:"locale"`
	Groups    []string `json:"groups"`
	Behaviour string   `json:"behaviour"`
	Value     []byte   `json:"value"`
}

type View struct {
	ID        string         `json:"id"`
	ProjectID string         `json:"projectId"`
	ShortID   string         `json:"shortId"`
	Name      string         `json:"name"`
	Variables []ViewVariable `json:"variables"`
}

type ViewVariable struct {
	ID      string `json:"id"`
	Locale  string `json:"locale"`
	ShortID string `json:"shortId"`
	Name    string `json:"name"`
}

type Model struct {
	Variables []VariableModel
	Name      string
	ProjectID string
	Locale    string
}

type LogicResult struct {
	ID        string
	ShortID   string
	ProjectID string
	Variables []ViewVariable
	Name      string
}

func NewModel(projectId, name string, entries []VariableModel) Model {
	return Model{
		Name:      name,
		ProjectID: projectId,
		Variables: entries,
	}
}

func (a *Model) Validate() map[string]string {
	v := map[string]interface{}{
		"projectID":          a.ProjectID,
		"groups":             nil,
		"name":               a.Name,
		"uniqueName":         a.Name,
		"validLocales":       nil,
		"validNum":           a.Variables,
		"validVariableNames": a.Variables,
		"behaviourValid":     a.Variables,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("name", validation.Required, validation.RuneLength(1, 200)),
			validation.Key("projectID", validation.Required, validation.RuneLength(26, 26)),
			validation.Key("uniqueName", validation.By(func(value interface{}) error {
				name := value.(string)

				var model declarations.Map
				res := storage.Gorm().Where("name = ? AND project_id = ?", name, a.ProjectID).Select("ID").First(&model)

				if errors.Is(res.Error, gorm.ErrRecordNotFound) {
					return nil
				}

				if res.Error != nil {
					return errors.New(fmt.Sprintf("Record with name '%s' already exists", name))
				}

				if model.ID != "" {
					return errors.New(fmt.Sprintf("Record with name '%s' already exists", name))
				}

				return nil
			})),
			validation.Key("validLocales", validation.By(func(value interface{}) error {
				return nil
			})),
			validation.Key("validNum", validation.By(func(value interface{}) error {
				if len(a.Variables) > 1000 {
					return errors.New("Number of map values cannot be larger than 1000.")
				}

				return nil
			})),
			validation.Key("validVariableNames", validation.By(func(value interface{}) error {
				m := make(map[string]int)
				for _, entry := range a.Variables {
					m[entry.Name] = 0
				}

				if len(m) != len(a.Variables) {
					return errors.New("Some variable/map names are not unique. All variable/map names must be unique.")
				}

				return nil
			})),
			validation.Key("behaviourValid", validation.By(func(value interface{}) error {
				if len(a.Variables) == 0 {
					return nil
				}

				m := make(map[string]string)
				for _, entry := range a.Variables {
					m[entry.Name] = entry.Behaviour
				}

				for key, v := range m {
					if v != constants.ReadonlyBehaviour && v != constants.ModifiableBehaviour {
						return errors.New(fmt.Sprintf("Invalid value for behaviour in variable '%s'. Variable behaviour can be 'modifiable' or 'readonly'", key))
					}
				}

				return nil
			})),
			validation.Key("groups", validation.By(func(value interface{}) error {
				for _, entry := range a.Variables {
					if len(entry.Groups) > 20 {
						return errors.New(fmt.Sprintf("Invalid number of groups for '%s'. Maximum number of groups per variable is 20.", entry.Name))
					}

					for _, g := range entry.Groups {
						if len(g) > 100 {
							return errors.New(fmt.Sprintf("Invalid group length for '%s'. Maximum number of characters per groups is 100.", g))
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

func newView(model LogicResult) View {
	return View{
		ID:        model.ID,
		ShortID:   model.ShortID,
		ProjectID: model.ProjectID,
		Name:      model.Name,
		Variables: model.Variables,
	}
}
