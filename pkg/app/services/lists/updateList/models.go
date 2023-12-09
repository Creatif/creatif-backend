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
	ID        string
	ShortID   string
	Values    ModelValues
	ProjectID string
}

func NewModel(projectId string, fields []string, name, id, shortID, updatingName string) Model {
	return Model{
		Fields:    fields,
		ProjectID: projectId,
		Name:      name,
		ID:        id,
		ShortID:   shortID,
		Values: ModelValues{
			Name: updatingName,
		},
	}
}

func (a *Model) Validate() map[string]string {
	v := map[string]interface{}{
		"fieldsValid":        a.Fields,
		"name":               a.Name,
		"id":                 a.ID,
		"idExists":           nil,
		"projectID":          a.ProjectID,
		"updatingNameExists": a.Values.Name,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("name", validation.When(a.Name != "", validation.RuneLength(1, 200))),
			validation.Key("id", validation.When(a.ID != "", validation.RuneLength(26, 26))),
			validation.Key("idExists", validation.By(func(value interface{}) error {
				name := a.Name
				shortId := a.ShortID
				id := a.ID

				if id != "" && len(id) != 26 {
					return errors.New("ID must have 26 characters")
				}

				if name == "" && shortId == "" && id == "" {
					return errors.New("At least one of 'id', 'name' or 'shortID' must be supplied in order to identify this list.")
				}
				return nil
			})),
			validation.Key("projectID", validation.Required, validation.RuneLength(26, 26)),
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
				res := storage.Gorm().Where("project_id = ? AND name = ?", a.ProjectID, a.Values.Name).Select("ID").First(&exists)
				if res.Error != nil && !errors.Is(res.Error, gorm.ErrRecordNotFound) {
					return errors.New(fmt.Sprintf("List with name '%s' already exists.", t))
				}

				if exists.ID != "" {
					return errors.New(fmt.Sprintf("List with name '%s' already exists.", t))
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
		ShortID:   model.ShortID,
		ProjectID: model.ProjectID,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}
