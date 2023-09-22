package createProject

import (
	"creatif/pkg/app/domain/app"
	"creatif/pkg/lib/sdk"
	"creatif/pkg/lib/storage"
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gorm.io/gorm"
	"time"
)

type Model struct {
	Name string `json:"name"`
}

func NewModel(name string) Model {
	return Model{
		Name: name,
	}
}

func (a *Model) Validate() map[string]string {
	v := map[string]interface{}{
		"name": a.Name,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("name", validation.Required, validation.RuneLength(1, 200), validation.By(func(value interface{}) error {
				name := value.(string)

				var project app.Project
				if err := storage.GetBy((app.Project{}).TableName(), "name", name, &project, "id"); err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
					return errors.New(fmt.Sprintf("Project with name '%s' already exists", name))
				}

				if project.ID != "" {
					return errors.New(fmt.Sprintf("Project with name '%s' already exists", name))
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
	ID   string `json:"id"`
	Name string `json:"name"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func newView(model app.Project) View {
	return View{
		ID:        model.ID,
		Name:      model.Name,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}
