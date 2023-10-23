package updateList

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/locales"
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
	ProjectID string
	Locale    string
}

func NewModel(projectId, locale string, fields []string, name, updatingName string) Model {
	return Model{
		Fields:    fields,
		ProjectID: projectId,
		Locale:    locale,
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
		"projectID":          a.ProjectID,
		"locale":             a.Locale,
		"updatingNameExists": a.Values.Name,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("name", validation.Required, validation.RuneLength(1, 200)),
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

				localeID, err := locales.GetIDWithAlpha(a.Locale)
				if err != nil {
					return errors.New(fmt.Sprintf("Locale '%s' not found.", a.Locale))
				}

				var exists declarations.List
				res := storage.Gorm().Where("project_id = ? AND name = ? AND locale_id = ?", a.ProjectID, a.Values.Name, localeID).Select("ID").First(&exists)
				if res.Error != nil && !errors.Is(res.Error, gorm.ErrRecordNotFound) {
					return errors.New(fmt.Sprintf("List with name '%s' already exists.", t))
				}

				if exists.ID != "" {
					return errors.New(fmt.Sprintf("List with name '%s' already exists.", t))
				}

				return nil
			})),
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

type View struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	ShortID   string `json:"shortID"`
	Locale    string `json:"locale"`
	ProjectID string `json:"projectID"`

	CreatedAt time.Time `gorm:"<-:create" json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func newView(model declarations.List, locale string) View {
	return View{
		ID:        model.ID,
		Name:      model.Name,
		ShortID:   model.ShortID,
		ProjectID: model.ProjectID,
		Locale:    locale,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}
