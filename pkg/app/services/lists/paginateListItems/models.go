package paginateListItems

import (
	"creatif/pkg/lib/sdk"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Model struct {
	ProjectID string
	ListName  string

	Limit          int
	Page           int
	Search         string
	Filters        map[string]string
	OrderBy        string
	OrderDirection string
	Groups         []string
}

func NewModel(projectId, listName, orderBy, direction string, limit, page int, groups []string, filters map[string]string) Model {
	return Model{
		ProjectID:      projectId,
		ListName:       listName,
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
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("projectId", validation.Required, validation.RuneLength(1, 26)),
		),
	); err != nil {
		return sdk.ErrorToResponseError(err)
	}

	return nil
}
