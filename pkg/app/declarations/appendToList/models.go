package appendToList

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/sdk"
	"errors"
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
	ProjectID string
	Variables []Variable
}

func NewModel(projectId, name string, variables []Variable) Model {
	return Model{
		Name:      name,
		ProjectID: projectId,
		Variables: variables,
	}
}

func (a Model) Validate() map[string]string {
	v := map[string]interface{}{
		"variableLen": len(a.Variables),
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("variableLen", validation.By(func(value interface{}) error {
				l := value.(int)

				if l > 1000 {
					return errors.New("The number of variables when creating a list cannot be higher than 1000.")
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

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func newView(model declarations.List) View {
	return View{
		ID:        model.ID,
		ProjectID: model.ProjectID,
		Name:      model.Name,

		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}
