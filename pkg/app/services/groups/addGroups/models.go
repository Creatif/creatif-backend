package addGroups

import (
	"creatif/pkg/lib/sdk"
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type GroupModel struct {
	ID     string
	Name   string
	Type   string
	Action string
}

type Model struct {
	ProjectID string
	Groups    []GroupModel
}

func NewModel(projectId string, groups []GroupModel) Model {
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
			validation.Key("projectID", validation.Required, validation.RuneLength(27, 27)),
			validation.Key("groups", validation.By(func(value interface{}) error {
				if a.Groups != nil {
					if len(a.Groups) > 200 {
						return errors.New("Maximum number of groups is 200.")
					}

					for _, g := range a.Groups {
						if len(g.Name) > 200 {
							return errors.New("Group name must have less than 200 characters.")
						}
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
