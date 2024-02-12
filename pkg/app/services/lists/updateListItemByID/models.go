package updateListItemByID

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/locales"
	"creatif/pkg/app/services/shared"
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
	"references",
	"locale",
}

type ModelValues struct {
	Name      string
	Metadata  []byte
	Groups    []string
	Locale    string
	Behaviour string
	Value     []byte
}

type Model struct {
	Fields     []string
	ListName   string
	ItemID     string
	Values     ModelValues
	ProjectID  string
	References []shared.UpdateReference
}

type LogicResult struct {
	Variable declarations.ListVariable
	Groups   []string
}

func NewModel(projectId, locale string, fields []string, listName, itemId, updatingName, behaviour string, groups []string, metadata, value []byte, references []shared.UpdateReference) Model {
	return Model{
		Fields:     fields,
		ProjectID:  projectId,
		ListName:   listName,
		ItemID:     itemId,
		References: references,
		Values: ModelValues{
			Name:      updatingName,
			Metadata:  metadata,
			Locale:    locale,
			Groups:    groups,
			Behaviour: behaviour,
			Value:     value,
		},
	}
}

func (a *Model) Validate() map[string]string {
	v := map[string]interface{}{
		"fieldsValid": a.Fields,
		"name":        a.ListName,
		"itemID":      a.ItemID,
		"projectID":   a.ProjectID,
		"locale":      a.Values.Locale,
		"groups":      a.Values.Groups,
		"behaviour":   a.Values.Behaviour,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("name", validation.Required),
			validation.Key("itemID", validation.Required),
			validation.Key("projectID", validation.Required, validation.RuneLength(26, 26)),
			validation.Key("fieldsValid", validation.Required, validation.By(func(value interface{}) error {
				t := value.([]string)

				if len(t) == 0 || len(t) > 7 {
					return errors.New(fmt.Sprintf("Invalid updateable fields. Valid updatable fields are %s", strings.Join(validUpdateableFields, ", ")))
				}

				if !sdk.ArrEqual(t, validUpdateableFields) {
					return errors.New(fmt.Sprintf("Invalid updateable fields. Valid updatable fields are %s", strings.Join(validUpdateableFields, ", ")))
				}

				return nil
			})),
			validation.Key("groups", validation.When(len(a.Values.Groups) != 0, validation.Each(validation.RuneLength(1, 200))), validation.By(func(value interface{}) error {
				groups := value.([]string)
				if len(groups) > 20 {
					return errors.New("Maximum number of groups is 20.")
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
			validation.Key("locale", validation.By(func(value interface{}) error {
				if !sdk.Includes(a.Fields, "locale") {
					return nil
				}

				t := value.(string)

				if !sdk.Includes(a.Fields, "locale") {
					return nil
				}

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
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	Locale    string      `json:"locale"`
	Groups    []string    `json:"groups"`
	Behaviour string      `json:"behaviour"`
	Metadata  interface{} `json:"metadata"`
	Value     interface{} `json:"value"`

	CreatedAt time.Time `gorm:"<-:createProject" json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func newView(model LogicResult) View {
	var m interface{} = model.Variable.Metadata
	if len(model.Variable.Metadata) == 0 {
		m = nil
	}

	var v interface{} = model.Variable.Value
	if len(model.Variable.Value) == 0 {
		v = nil
	}

	locale, _ := locales.GetAlphaWithID(model.Variable.LocaleID)
	return View{
		ID:        model.Variable.ID,
		Locale:    locale,
		Name:      model.Variable.Name,
		Groups:    model.Groups,
		Behaviour: model.Variable.Behaviour,
		Metadata:  m,
		Value:     v,
		CreatedAt: model.Variable.CreatedAt,
		UpdatedAt: model.Variable.UpdatedAt,
	}
}
