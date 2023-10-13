package removeMapEntry

import (
	"creatif/pkg/app/services/locales"
	"creatif/pkg/lib/sdk"
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Model struct {
	Name      string `json:"name"`
	EntryName string `json:"entryName"`
	ProjectID string `json:"projectID"`
	Locale    string `json:"locale"`
}

func NewModel(projectId, locale, name, entryName string) Model {
	return Model{
		Name:      name,
		Locale:    locale,
		ProjectID: projectId,
		EntryName: entryName,
	}
}

func (a Model) Validate() map[string]string {
	v := map[string]interface{}{
		"projectID": a.ProjectID,
		"entryName": a.ProjectID,
		"locale":    a.Locale,
		"name":      a.Name,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("name", validation.Required, validation.RuneLength(1, 200)),
			validation.Key("entryName", validation.Required, validation.RuneLength(1, 200)),
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
