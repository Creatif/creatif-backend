package addGroups

import (
	"creatif/pkg/lib/sdk"
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Model struct {
	ProjectID string
	Groups    []string
}

func NewModel(projectId string, groups []string) Model {
	return Model{
		ProjectID: projectId,
		Groups:    groups,
	}
}

func (a *Model) Validate() map[string]string {
	v := map[string]interface{}{
		"groups":    a.Groups,
		"projectID": a.ProjectID,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("projectID", validation.Required, validation.RuneLength(26, 26)),
			validation.Key("groups", validation.When(len(a.Groups) != 0, validation.Each(validation.RuneLength(1, 200))), validation.By(func(value interface{}) error {
				if a.Groups != nil {
					if len(a.Groups) > 200 {
						return errors.New("Maximum number of groups is 200.")
					}

					return nil
				}

				// add unique group check

				return nil
			})),
		),
	); err != nil {
		return sdk.ErrorToResponseError(err)
	}

	return nil
}
