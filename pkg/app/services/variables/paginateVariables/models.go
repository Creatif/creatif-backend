package paginateVariables

import (
	"creatif/pkg/app/services/locales"
	"creatif/pkg/lib/sdk"
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Model struct {
	ProjectID string
	Locale    string

	Limit          int
	Page           int
	Search         string
	Filters        map[string]string
	OrderBy        string
	OrderDirection string
	Groups         []string
}

func NewModel(projectId, locale, orderBy, direction string, limit, page int, groups []string, filters map[string]string) Model {
	return Model{
		ProjectID:      projectId,
		Locale:         locale,
		OrderBy:        orderBy,
		Page:           page,
		Filters:        filters,
		OrderDirection: direction,
		Limit:          limit,
		Groups:         groups,
	}
}

func (a *Model) Validate() map[string]string {
	v := map[string]interface{}{
		"projectId": a.ProjectID,
		"locale":    a.Locale,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("projectId", validation.Required, validation.RuneLength(26, 26)),
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
