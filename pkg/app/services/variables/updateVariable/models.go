package updateVariable

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/locales"
	"creatif/pkg/lib/constants"
	"creatif/pkg/lib/sdk"
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"strings"
	"time"
)

var validUpdateableFields = []string{
	"name",
	"metadata",
	"groups",
	"behaviour",
	"value",
	"locale",
}

type ModelValues struct {
	Name      string
	Metadata  []byte
	Groups    []string
	Behaviour string
	Locale    string
	Value     []byte
}

type Model struct {
	Fields    []string
	ID        string
	Values    ModelValues
	ProjectID string
}

func NewModel(projectId string, fields []string, name, updatingName, behaviour string, groups []string, metadata, value []byte, updatingLocale string) Model {
	return Model{
		Fields:    fields,
		ProjectID: projectId,
		ID:        name,
		Values: ModelValues{
			Name:      updatingName,
			Metadata:  metadata,
			Groups:    groups,
			Locale:    updatingLocale,
			Behaviour: behaviour,
			Value:     value,
		},
	}
}

func (a *Model) Validate() map[string]string {
	v := map[string]interface{}{
		"projectID":       a.ProjectID,
		"updatingLocale":  a.Values.Locale,
		"fieldsValid":     a.Fields,
		"updatingName":    a.Values.Name,
		"ID":              a.ID,
		"behaviour":       a.Values.Behaviour,
		"groups":          a.Values.Groups,
		"nameLocaleValid": nil,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("projectID", validation.Required, validation.RuneLength(26, 26)),
			validation.Key("ID", validation.Required),
			validation.Key("updatingName", validation.When(a.Values.Name != "", validation.RuneLength(1, 200))),
			validation.Key("updatingLocale", validation.By(func(value interface{}) error {
				if !sdk.Includes(a.Fields, "locale") {
					return nil
				}

				t := value.(string)

				if !locales.ExistsByAlpha(t) {
					return errors.New(fmt.Sprintf("Locale '%s' does not exist.", t))
				}

				return nil
			})),
			validation.Key("fieldsValid", validation.Required, validation.By(func(value interface{}) error {
				t := value.([]string)

				if len(t) == 0 || len(t) > 6 {
					return errors.New(fmt.Sprintf("Invalid updateable fields. Valid updatable fields are %s", strings.Join(validUpdateableFields, ", ")))
				}

				if !sdk.ArrEqual(t, validUpdateableFields) {
					return errors.New(fmt.Sprintf("Invalid updateable fields. Valid updatable fields are %s", strings.Join(validUpdateableFields, ", ")))
				}

				return nil
			})),
			validation.Key("nameLocaleValid", validation.By(func(value interface{}) error {
				if sdk.Includes(a.Fields, "locale") || sdk.Includes(a.Fields, "name") {
					if a.Values.Name == "" || a.Values.Locale == "" {
						return errors.New("If updating either 'name' or 'locale', both name and locale must exist as updating values.")
					}
				}

				return nil
			})),
			validation.Key("behaviour", validation.By(func(value interface{}) error {
				if !sdk.Includes(a.Fields, "behaviour") {
					return nil
				}

				t := value.(string)

				if t != constants.ReadonlyBehaviour && t != constants.ModifiableBehaviour {
					return errors.New(fmt.Sprintf("Invalid value for behaviour. Variable behaviour can be 'modifiable' or 'readonly'"))
				}

				return nil
			})),
			validation.Key("groups", validation.When(len(a.Values.Groups) != 0, validation.Each(validation.RuneLength(1, 100))), validation.By(func(value interface{}) error {
				groups := value.([]string)
				if len(groups) > 20 {
					return errors.New(fmt.Sprintf("Invalid number of groups for '%s'. Maximum number of groups per variable is 20.", a.ID))
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
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	Groups    []string    `json:"groups"`
	ShortID   string      `json:"shortID"`
	Behaviour string      `json:"behaviour"`
	Metadata  interface{} `json:"metadata"`
	Locale    string      `json:"locale"`
	Value     interface{} `json:"value"`

	CreatedAt time.Time `gorm:"<-:createProject" json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func newView(model declarations.Variable) View {
	var m interface{} = model.Metadata
	if len(model.Metadata) == 0 {
		m = nil
	}

	var v interface{} = model.Value
	if len(model.Value) == 0 {
		v = nil
	}

	locale, _ := locales.GetAlphaWithID(model.LocaleID)

	return View{
		ID:        model.ID,
		Name:      model.Name,
		ShortID:   model.ShortID,
		Locale:    locale,
		Groups:    model.Groups,
		Behaviour: model.Behaviour,
		Metadata:  m,
		Value:     v,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}
