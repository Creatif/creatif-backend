package deleteListItemByIndex

import (
	"creatif/pkg/app/services/locales"
	"creatif/pkg/lib/sdk"
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Model struct {
	Name      string
	ItemIndex int64
	ProjectID string
	Locale    string
}

func NewModel(projectId, locale, name string, itemIndex int64) Model {
	return Model{
		Name:      name,
		ItemIndex: itemIndex,
		ProjectID: projectId,
		Locale:    locale,
	}
}

func (a Model) Validate() map[string]string {
	v := map[string]interface{}{
		"name":      a.Name,
		"projectID": a.ProjectID,
		"locale":    a.Locale,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("name", validation.Required, validation.RuneLength(1, 200)),
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
