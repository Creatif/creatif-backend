package updateList

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/sdk"
	"creatif/pkg/lib/storage"
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gorm.io/gorm"
	"strings"
	"time"
)

var validUpdateableFields = []string{
	"name",
}

type ModelValues struct {
	Name string `json:"name"`
}

type Model struct {
	Fields    []string
	Name      string
	Values    ModelValues
	ProjectID string `json:"projectID"`
}

func NewModel(projectId string, fields []string, name, updatingName string) Model {
	return Model{
		Fields:    fields,
		ProjectID: projectId,
		Name:      name,
		Values: ModelValues{
			Name: updatingName,
		},
	}
}

func (a *Model) Validate() map[string]string {
	v := map[string]interface{}{
		"fieldsValid":        a.Fields,
		"name":               a.Values.Name,
		"updatingNameExists": a.Values.Name,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("name", validation.Required, validation.RuneLength(1, 200)),
			validation.Key("fieldsValid", validation.Required, validation.By(func(value interface{}) error {
				t := value.([]string)

				if len(t) == 0 || len(t) > 1 {
					return errors.New(fmt.Sprintf("Invalid updateable fields. Valid updatable fields are '%s'", strings.Join(validUpdateableFields, ", ")))
				}

				if !sdk.ArrEqual(t, validUpdateableFields) {
					return errors.New(fmt.Sprintf("Invalid updateable fields. Valid updatable fields are '%s'", strings.Join(validUpdateableFields, ", ")))
				}

				return nil
			})),
			validation.Key("updatingNameExists", validation.When(a.Values.Name != "", validation.Required, validation.RuneLength(1, 200)), validation.By(func(value interface{}) error {
				if !sdk.Includes(a.Fields, "name") {
					return nil
				}

				t := value.(string)

				if t == "" {
					return nil
				}

				var exists declarations.List
				if err := storage.GetBy((declarations.List{}).TableName(), "name", t, &exists, "id"); !errors.Is(err, gorm.ErrRecordNotFound) {
					return errors.New(fmt.Sprintf("Variable with name '%s' already exists.", t))
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
	Name      string `json:"name"`
	ShortID   string `json:"shortID"`
	ProjectID string `json:"projectID"`

	CreatedAt time.Time `gorm:"<-:create" json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func newView(model declarations.List) View {
	return View{
		ID:        model.ID,
		Name:      model.Name,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}
