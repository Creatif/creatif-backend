package replaceListItem

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/sdk"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"time"
)

type Variable struct {
	Name      string
	Metadata  []byte
	Groups    []string
	Behaviour string
	Value     []byte
}

type Model struct {
	Name      string
	ItemName  string
	ProjectID string
	Variable  Variable
}

func NewModel(projectId, name, itemName string, variables Variable) Model {
	return Model{
		Name:      name,
		ItemName:  itemName,
		ProjectID: projectId,
		Variable:  variables,
	}
}

func (a Model) Validate() map[string]string {
	v := map[string]interface{}{
		"name":     a.Name,
		"itemName": a.ItemName,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("name", validation.Required, validation.RuneLength(1, 200)),
			validation.Key("itemName", validation.Required, validation.RuneLength(1, 200)),
		),
	); err != nil {
		return sdk.ErrorToResponseError(err)
	}

	return nil
}

type View struct {
	ID        string
	Index     string
	Name      string
	Metadata  []byte
	Groups    []string
	Behaviour string
	Value     []byte
	CreatedAt time.Time
	UpdatedAt time.Time
}

func newView(model declarations.ListVariable) View {
	return View{}
}
