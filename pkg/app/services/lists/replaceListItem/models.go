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
	Name        string
	ID          string
	ShortID     string
	ItemID      string
	ItemShortID string
	ProjectID   string
	Locale      string
	Variable    Variable
}

func NewModel(projectId, locale, name, id, shortID, itemID, itemShortID string, variable Variable) Model {
	return Model{
		Name:        name,
		ID:          id,
		ShortID:     shortID,
		Locale:      locale,
		ItemID:      itemID,
		ItemShortID: itemShortID,
		ProjectID:   projectId,
		Variable:    variable,
	}
}

func (a Model) Validate() map[string]string {
	v := map[string]interface{}{
		"name":         a.Name,
		"id":           a.ID,
		"idExists":     nil,
		"projectID":    a.ProjectID,
		"locale":       a.Locale,
		"itemIDExists": nil,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("name", validation.When(a.Name != "", validation.RuneLength(1, 200))),
			validation.Key("id", validation.When(a.ID != "", validation.RuneLength(26, 26))),
			validation.Key("idExists", validation.By(func(value interface{}) error {
				name := a.Name
				shortId := a.ShortID
				id := a.ID

				if id != "" && len(id) != 26 {
					return errors.New("ID must have 26 characters")
				}

				if name == "" && shortId == "" && id == "" {
					return errors.New("At least one of 'id', 'name' or 'shortID' must be supplied in order to identify this list.")
				}
				return nil
			})),
			validation.Key("projectID", validation.Required, validation.RuneLength(26, 26)),
			validation.Key("itemIDExists", validation.By(func(value interface{}) error {
				id := a.ItemID
				shortId := a.ItemShortID

				if id != "" && len(id) != 26 {
					return errors.New("ID must have 26 characters")
				}

				if shortId == "" && id == "" {
					return errors.New("At least one of 'id' or 'shortID' must be supplied in order to identify this variable.")
				}
				return nil
			})),
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
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	Metadata  interface{} `json:"metadata"`
	Groups    []string    `json:"groups"`
	Behaviour string      `json:"behaviour"`
	Value     interface{} `json:"value"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"UpdatedAt"`
}

func newView(model declarations.ListVariable) View {
	return View{
		ID:        model.ID,
		Name:      model.Name,
		Metadata:  model.Metadata,
		Groups:    model.Groups,
		Behaviour: model.Behaviour,
		Value:     model.Value,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}
