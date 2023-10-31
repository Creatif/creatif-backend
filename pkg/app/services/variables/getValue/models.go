package getValue

import (
	"creatif/pkg/app/services/locales"
	"creatif/pkg/lib/sdk"
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gorm.io/datatypes"
)

type Model struct {
	Name      string `json:"name"`
	ID        string `json:"id"`
	ShortID   string `json:"shortID"`
	ProjectID string `json:"projectID"`
	Locale    string `json:"locale"`
}

func NewModel(projectId, id, shortId, name, locale string) Model {
	return Model{Name: name, ProjectID: projectId, Locale: locale, ID: id, ShortID: shortId}
}

type Variable struct {
	Value datatypes.JSON
}

func newView(model Variable) datatypes.JSON {
	return model.Value
}

func (a *Model) Validate() map[string]string {
	v := map[string]interface{}{
		"name":      a.Name,
		"id":        a.ID,
		"idExists":  nil,
		"projectID": a.ProjectID,
		"locale":    a.Locale,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("name", validation.When(a.Name != "", validation.RuneLength(1, 200))),
			validation.Key("id", validation.When(a.ID != "", validation.RuneLength(26, 26))),
			validation.Key("idExists", validation.By(func(value interface{}) error {
				name := a.Name
				shortId := a.ShortID
				id := a.ID

				if name == "" && shortId == "" && id == "" {
					return errors.New("At least one of 'id', 'name' or 'shortID' must be supplied in order to identify this variable.")
				}
				return nil
			})),
			validation.Key("projectID", validation.Required, validation.RuneLength(26, 26)),
			validation.Key("locale", validation.Required, validation.By(func(value interface{}) error {
				t := value.(string)

				if !locales.ExistsByAlpha(t) {
					return errors.New(fmt.Sprintf("Locale '%s' does not exist.", t))
				}

				return nil
			})),
		),
	); err != nil {
		return sdk.ErrorToResponseError(err)
	}

	return nil
}
