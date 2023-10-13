package mapCreate

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/locales"
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
	ProjectID string              `json:"projectID"`
	Locale    string              `json:"locale"`
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
	Locale    string  `json:"locale"`
}

type LogicResult struct {
	ID        string
	Locale    string
	ProjectID string
	Variables []map[string]string
	Name      string
}

func NewModel(projectId, locale, name string, entries []Entry) Model {
	return Model{
		Name:      name,
		ProjectID: projectId,
		Locale:    locale,
		Entries:   entries,
	}
}

func (a *Model) Validate() map[string]string {
	v := map[string]interface{}{
		"projectID":          a.ProjectID,
		"locale":             a.Locale,
		"groups":             nil,
		"name":               a.Name,
		"uniqueName":         a.Name,
		"validNum":           a.Entries,
		"validVariableNames": a.Entries,
		"behaviourValid":     a.Entries,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("name", validation.Required, validation.RuneLength(1, 200)),
			validation.Key("projectID", validation.Required, validation.RuneLength(26, 26)),
			validation.Key("locale", validation.Required, validation.By(func(value interface{}) error {
				t := value.(string)

				if !locales.ExistsByAlpha(t) {
					return errors.New(fmt.Sprintf("Locale '%s' does not exist.", t))
				}

				return nil
			})),
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
			validation.Key("groups", validation.By(func(value interface{}) error {
				for _, entry := range a.Entries {
					if entry.Type == "variable" {
						o := entry.Model.(VariableModel)

						if len(o.Groups) > 20 {
							return errors.New(fmt.Sprintf("Invalid number of groups for '%s'. Maximum number of groups per variable is 20.", o.Name))
						}

						for _, g := range o.Groups {
							if len(g) > 200 {
								return errors.New(fmt.Sprintf("Invalid group length for '%s'. Maximum number of characters per groups is 200.", g))
							}
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
		Locale:    model.Locale,
		ProjectID: model.ProjectID,
		Name:      model.Name,
		Variables: model.Variables,
	}
}
