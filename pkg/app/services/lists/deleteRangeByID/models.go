package deleteRangeByID

import (
	"creatif/pkg/app/services/locales"
	"creatif/pkg/lib/sdk"
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Model struct {
	Name      string
	Items     []string
	ProjectID string
	Locale    string
}

func NewModel(projectId, locale, name string, items []string) Model {
	return Model{
		Name:      name,
		Locale:    locale,
		Items:     items,
		ProjectID: projectId,
	}
}

func (a Model) Validate() map[string]string {
	v := map[string]interface{}{
		"name":      a.Name,
		"projectID": a.ProjectID,
		"items":     a.Items,
		"locale":    a.Locale,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("name", validation.Required, validation.RuneLength(1, 200)),
			validation.Key("projectID", validation.Required, validation.RuneLength(26, 26)),
			validation.Key("items", validation.Required, validation.By(func(value interface{}) error {
				items := a.Items
				if items == nil {
					return errors.New("Items number must be bigger than 0 (zero).")
				}

				if len(items) == 0 {
					return errors.New("Items number must be bigger than 0 (zero).")
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
