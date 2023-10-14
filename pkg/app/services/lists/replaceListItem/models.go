package replaceListItem

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
	Value     []byte
}

type Model struct {
	Name      string
	ItemName  string
	ProjectID string
	Locale    string
	Variable  Variable
}

func NewModel(projectId, locale, name, itemName string, variable Variable) Model {
	return Model{
		Name:      name,
		Locale:    locale,
		ItemName:  itemName,
		ProjectID: projectId,
		Variable:  variable,
	}
}

func (a Model) Validate() map[string]string {
	v := map[string]interface{}{
		"name":      a.Name,
		"projectID": a.ProjectID,
		"locale":    a.Locale,
		"itemName":  a.ItemName,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("projectID", validation.Required, validation.RuneLength(26, 26)),
			validation.Key("name", validation.Required, validation.RuneLength(1, 200)),
			validation.Key("itemName", validation.Required, validation.RuneLength(1, 200)),
			validation.Key("locale", validation.Required, validation.By(func(value interface{}) error {
				t := value.(string)

				if !locales.ExistsByAlpha(t) {
					return errors.New(fmt.Sprintf("Locale '%s' not found.", t))
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
	return View{
		ID:        model.ID,
		Index:     model.Index,
		Name:      model.Name,
		Metadata:  model.Metadata,
		Groups:    model.Groups,
		Behaviour: model.Behaviour,
		Value:     model.Value,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}
