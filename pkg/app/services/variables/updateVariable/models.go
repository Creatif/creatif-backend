package updateVariable

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/locales"
	"creatif/pkg/lib/constants"
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
	"metadata",
	"groups",
	"behaviour",
	"value",
}

type ModelValues struct {
	Name      string
	Metadata  []byte
	Groups    []string
	Behaviour string
	Value     []byte
}

type Model struct {
	Fields    []string
	Name      string
	ID        string
	ShortID   string
	Values    ModelValues
	ProjectID string
	Locale    string
}

func NewModel(projectId, locale string, fields []string, name, id, shortId, updatingName, behaviour string, groups []string, metadata, value []byte) Model {
	return Model{
		Fields:    fields,
		ProjectID: projectId,
		Name:      name,
		ID:        id,
		ShortID:   shortId,
		Locale:    locale,
		Values: ModelValues{
			Name:      updatingName,
			Metadata:  metadata,
			Groups:    groups,
			Behaviour: behaviour,
			Value:     value,
		},
	}
}

func (a *Model) Validate() map[string]string {
	v := map[string]interface{}{
		"projectID":          a.ProjectID,
		"locale":             a.Locale,
		"fieldsValid":        a.Fields,
		"updatingName":       a.Values.Name,
		"name":               a.Name,
		"id":                 a.ID,
		"idExists":           nil,
		"behaviour":          a.Values.Behaviour,
		"updatingNameExists": a.Values.Name,
		"groups":             a.Values.Groups,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("projectID", validation.Required, validation.RuneLength(26, 26)),
			validation.Key("name", validation.When(a.Name != "", validation.RuneLength(1, 200))),
			validation.Key("updatingName", validation.When(a.Name != "", validation.RuneLength(1, 200))),
			validation.Key("id", validation.When(a.ID != "", validation.RuneLength(26, 26))),
			validation.Key("idExists", validation.By(func(value interface{}) error {
				name := a.Name
				shortId := a.ShortID
				id := a.ID

				if name == "" && shortId == "" && id == "" {
					return errors.New("At least one of 'id', 'name' or 'shortID' must be supplied in order to identify this variable.")
				}
				return nil
			})),
			validation.Key("locale", validation.Required, validation.By(func(value interface{}) error {
				t := value.(string)

				if !locales.ExistsByAlpha(t) {
					return errors.New(fmt.Sprintf("Locale '%s' does not exist.", t))
				}

				return nil
			})),
			validation.Key("fieldsValid", validation.Required, validation.By(func(value interface{}) error {
				t := value.([]string)

				if len(t) == 0 || len(t) > 5 {
					return errors.New(fmt.Sprintf("Invalid updateable fields. Valid updatable fields are %s", strings.Join(validUpdateableFields, ", ")))
				}

				if !sdk.ArrEqual(t, validUpdateableFields) {
					return errors.New(fmt.Sprintf("Invalid updateable fields. Valid updatable fields are %s", strings.Join(validUpdateableFields, ", ")))
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
			validation.Key("updatingNameExists", validation.When(a.Values.Name != "", validation.Required, validation.RuneLength(1, 200)), validation.By(func(value interface{}) error {
				if !sdk.Includes(a.Fields, "name") {
					return nil
				}

				t := value.(string)

				if t == "" {
					return nil
				}

				var exists declarations.Variable
				if err := storage.GetBy((declarations.Variable{}).TableName(), "name", t, &exists, "id"); !errors.Is(err, gorm.ErrRecordNotFound) {
					return errors.New(fmt.Sprintf("Variable with name '%s' already exists.", t))
				}

				return nil
			})),
			validation.Key("groups", validation.When(len(a.Values.Groups) != 0, validation.Each(validation.RuneLength(1, 100))), validation.By(func(value interface{}) error {
				groups := value.([]string)
				if len(groups) > 20 {
					return errors.New(fmt.Sprintf("Invalid number of groups for '%s'. Maximum number of groups per variable is 20.", a.Name))
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

func newView(model declarations.Variable, locale string) View {
	return View{
		ID:        model.ID,
		Name:      model.Name,
		ShortID:   model.ShortID,
		Locale:    locale,
		Groups:    model.Groups,
		Behaviour: model.Behaviour,
		Metadata:  model.Metadata,
		Value:     model.Value,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}
