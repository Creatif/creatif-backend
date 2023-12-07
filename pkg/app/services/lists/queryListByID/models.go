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
	Name      string
	ItemID    string
	ProjectID string
	Locale    string
}

func NewModel(projectId, locale, name, itemID string) Model {
	return Model{
		ProjectID: projectId,
		Locale:    locale,
		Name:      name,
		ItemID:    itemID,
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
		"name":      a.Name,
		"itemId":    a.ItemID,
		"projectID": a.ProjectID,
		"locale":    a.Locale,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("name", validation.Required),
			validation.Key("itemId", validation.Required),
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
