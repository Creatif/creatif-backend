package queryListByID

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/locales"
	"creatif/pkg/lib/sdk"
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/lib/pq"
	"time"
)

type Model struct {
	Name        string
	ID          string
	ShortID     string
	ItemID      string
	ItemShortID string
	ProjectID   string
	Locale      string
}

func NewModel(projectId, locale, name, id, shortID, itemID, itemShortID string) Model {
	return Model{
		ProjectID:   projectId,
		Locale:      locale,
		Name:        name,
		ShortID:     shortID,
		ID:          id,
		ItemID:      itemID,
		ItemShortID: itemShortID,
	}
}

type View struct {
	ID        string         `json:"id"`
	Locale    string         `json:"locale"`
	ShortID   string         `json:"shortId"`
	Name      string         `json:"name"`
	Behaviour string         `json:"behaviour"`
	Groups    pq.StringArray `json:"groups"`
	Metadata  interface{}    `json:"metadata"`
	Value     interface{}    `json:"value"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func newView(model declarations.ListVariable, locale string) View {
	return View{
		ID:        model.ID,
		Locale:    locale,
		ShortID:   model.ShortID,
		Name:      model.Name,
		Behaviour: model.Behaviour,
		Groups:    model.Groups,
		Metadata:  model.Metadata,
		Value:     model.Value,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}

func (a *Model) Validate() map[string]string {
	v := map[string]interface{}{
		"name":         a.Name,
		"id":           a.ID,
		"idExists":     nil,
		"itemIDExists": nil,
		"projectID":    a.ProjectID,
		"locale":       a.Locale,
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
			validation.Key("itemIDExists", validation.By(func(value interface{}) error {
				shortId := a.ItemShortID
				id := a.ItemID

				if id != "" && len(id) != 26 {
					return errors.New("ID must have 26 characters")
				}

				if shortId == "" && id == "" {
					return errors.New("At least one of 'id', or 'shortID' must be supplied in order to identify this variable.")
				}
				return nil
			})),
			validation.Key("projectID", validation.Required, validation.RuneLength(26, 26)),
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
