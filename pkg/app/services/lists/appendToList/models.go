package appendToList

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
	ID        string
	ShortID   string
	ProjectID string
	Locale    string
	Variables []Variable
}

func NewModel(projectId, locale, name, id, shortID string, variables []Variable) Model {
	return Model{
		Name:      name,
		ID:        id,
		ShortID:   shortID,
		Locale:    locale,
		ProjectID: projectId,
		Variables: variables,
	}
}

func (a Model) Validate() map[string]string {
	v := map[string]interface{}{
		"variableLen": len(a.Variables),
		"projectID":   a.ProjectID,
		"locale":      a.Locale,
		"name":        a.Name,
		"id":          a.ID,
		"idExists":    nil,
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
			validation.Key("variableLen", validation.By(func(value interface{}) error {
				l := value.(int)

				if l > 1000 {
					return errors.New("The number of variables when creating a list cannot be higher than 1000.")
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
	ID        string `json:"id"`
	ProjectID string `json:"projectID"`
	Name      string `json:"name"`
	Locale    string `json:"locale"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func newView(model declarations.List, locale string) View {
	return View{
		ID:        model.ID,
		Locale:    locale,
		ProjectID: model.ProjectID,
		Name:      model.Name,

		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}
