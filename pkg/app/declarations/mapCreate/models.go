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
	Type      string   `json:"type"`
	Metadata  []byte   `json:"metadata"`
	Groups    []string `json:"groups"`
	Behaviour string   `json:"behaviour"`
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
	Nodes []string
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
		"uniqueName": a.Name,
		"validNum":   a.Nodes,
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
				names := value.([]string)
				if len(names) > 100 {
					return errors.New("Number of nodes cannot be higher than 100")
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
	ID    string   `json:"id"`
	Name  string   `json:"name"`
	Nodes []string `json:"nodes"`
}

func newView(model LogicResult) View {
	return View{
		ID:    model.ID,
		Name:  model.Name,
		Nodes: model.Nodes,
	}
}
