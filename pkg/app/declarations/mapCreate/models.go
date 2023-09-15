package mapCreate

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/sdk"
	"creatif/pkg/lib/storage"
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gorm.io/gorm"
)

type NodeModel struct {
	Name      string   `json:"name"`
	Metadata  []byte   `json:"metadata"`
	Groups    []string `json:"groups"`
	Behaviour string   `json:"behaviour"`
	Value     []byte   `json:"value"`
}

type Entry struct {
	Type  string
	Model interface{}
}

type Model struct {
	Entries []Entry `json:"entries"`
	Name    string  `json:"name"`
}

type LogicResult struct {
	ID    string
	Nodes []map[string]string
	Name  string
}

func NewModel(name string, entries []Entry) Model {
	return Model{
		Name:    name,
		Entries: entries,
	}
}

func (a *Model) Validate() map[string]string {
	v := map[string]interface{}{
		"uniqueName":     a.Name,
		"validNum":       a.Entries,
		"validNodeNames": a.Entries,
	}

	if err := validation.Validate(v,
		validation.Map(
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
			validation.Key("validNodeNames", validation.By(func(value interface{}) error {
				m := make(map[string]int)
				for _, entry := range a.Entries {
					if entry.Type == "node" {
						o := entry.Model.(NodeModel)
						m[o.Name] = 0
					}
				}

				if len(m) != len(a.Entries) {
					return errors.New("Some node/map names are not unique. All node/map names must be unique.")
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
	ID    string              `json:"id"`
	Name  string              `json:"name"`
	Nodes []map[string]string `json:"nodes"`
}

func newView(model LogicResult) View {
	return View{
		ID:    model.ID,
		Name:  model.Name,
		Nodes: model.Nodes,
	}
}
