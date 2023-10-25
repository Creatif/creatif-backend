package updateMapVariable

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
}

type VariableModel struct {
	Name      string
	Metadata  []byte
	Groups    []string
	Behaviour string
	Value     []byte
}

type Model struct {
	Fields       []string
	Values       VariableModel
	MapName      string
	VariableName string
	ProjectID    string
	Locale       string
}

func NewModel(projectId, locale, mapName, variableName string, fields []string, values VariableModel) Model {
	return Model{
		MapName:      mapName,
		Locale:       locale,
		Fields:       fields,
		ProjectID:    projectId,
		Values:       values,
		VariableName: variableName,
	}
}

type LogicResult struct {
	Map    declarations.Map
	Entry  declarations.MapVariable
	Locale string
}

func (a *Model) Validate() map[string]string {
	v := map[string]interface{}{
		"groups":       a.Values.Groups,
		"mapName":      a.MapName,
		"fieldsValid":  a.Fields,
		"variableName": a.VariableName,
		"behaviour":    a.Values.Behaviour,
		"projectID":    a.ProjectID,
		"locale":       a.Locale,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("mapName", validation.Required, validation.RuneLength(1, 200)),
			validation.Key("variableName", validation.Required, validation.RuneLength(1, 200)),
			validation.Key("projectID", validation.Required, validation.RuneLength(26, 26)),
			validation.Key("behaviour", validation.Required, validation.By(func(value interface{}) error {
				v := value.(string)
				if v != constants.ReadonlyBehaviour && v != constants.ModifiableBehaviour {
					return errors.New(fmt.Sprintf("Invalid value for behaviour in variable '%s'. Variable behaviour can be 'modifiable' or 'readonly'", v))
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
			validation.Key("groups", validation.When(len(a.Values.Groups) != 0, validation.Each(validation.RuneLength(1, 100))), validation.By(func(value interface{}) error {
				if a.Values.Groups != nil {
					if len(a.Values.Groups) > 20 {
						return errors.New("Maximum number of groups is 20.")
					}

					return nil
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
		),
	); err != nil {
		return sdk.ErrorToResponseError(err)
	}

	return nil
}

type ViewEntry struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	ShortID   string      `json:"shortID"`
	Metadata  interface{} `json:"metadata"`
	Groups    []string    `json:"groups"`
	Behaviour string      `json:"behaviour"`
	Value     interface{} `json:"value"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
}

type View struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	ProjectID string `json:"projectID"`
	Locale    string `json:"locale"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Entry ViewEntry `json:"entry"`
}

func newView(logicResult LogicResult) View {
	m := logicResult.Map
	variable := logicResult.Entry

	return View{
		ID:        m.ID,
		Name:      m.Name,
		Locale:    logicResult.Locale,
		ProjectID: m.ProjectID,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		Entry: ViewEntry{
			ID:        variable.ID,
			Name:      variable.Name,
			ShortID:   variable.ShortID,
			Metadata:  variable.Metadata,
			Groups:    variable.Groups,
			Behaviour: variable.Behaviour,
			Value:     variable.Value,
			CreatedAt: variable.CreatedAt,
			UpdatedAt: variable.UpdatedAt,
		},
	}
}
