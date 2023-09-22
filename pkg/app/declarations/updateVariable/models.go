package updateVariable

import (
	"creatif/pkg/app/domain/declarations"
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
	Name      string   `json:"name"`
	Metadata  []byte   `json:"metadata"`
	Groups    []string `json:"groups"`
	Behaviour string   `json:"behaviour"`
	Value     []byte   `json:"value"`
}

type Model struct {
	Fields []string
	Name   string
	Values ModelValues
}

func NewModel(fields []string, name, updatingName, behaviour string, groups []string, metadata, value []byte) Model {
	return Model{
		Fields: fields,
		Name:   name,
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
		"fieldsValid":        a.Fields,
		"name":               a.Values.Name,
		"groups":             a.Values.Groups,
		"behaviour":          a.Values.Behaviour,
		"updatingNameExists": a.Values.Name,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("name", validation.Required, validation.RuneLength(1, 200)),
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
			validation.Key("groups", validation.When(len(a.Values.Groups) != 0, validation.Each(validation.RuneLength(1, 200)))),
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
	Behaviour string      `json:"behaviour"`
	Metadata  interface{} `json:"metadata"`
	Value     interface{} `json:"value"`

	CreatedAt time.Time `gorm:"<-:createVariable" json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func newView(model declarations.Variable) View {
	return View{
		ID:        model.ID,
		Name:      model.Name,
		Groups:    model.Groups,
		Behaviour: model.Behaviour,
		Metadata:  model.Metadata,
		Value:     model.Value,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}