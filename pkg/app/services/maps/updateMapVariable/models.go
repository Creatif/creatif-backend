package updateMapVariable

import (
	"creatif/pkg/app/services/locales"
	"creatif/pkg/lib/constants"
	"creatif/pkg/lib/sdk"
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/lib/pq"
	"gorm.io/datatypes"
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

type MapVariableWithMap struct {
	ID      string `gorm:"primarykey;type:text;default:gen_ulid()"`
	ShortID string `gorm:"uniqueIndex:unique_variable;type:text"`

	Name      string `gorm:"uniqueIndex:unique_map_variable"`
	Behaviour string
	Groups    pq.StringArray `gorm:"type:text[]"`
	Metadata  datatypes.JSON `gorm:"type:jsonb"`
	Value     datatypes.JSON `gorm:"type:jsonb"`

	MapID    string `gorm:"column:map_id"`
	LocaleID string `gorm:"uniqueIndex:unique_map_variable;type:text"`

	CreatedAt time.Time
	UpdatedAt time.Time

	MapName      string    `gorm:"column:map_name"`
	MapCreatedAt time.Time `gorm:"column:map_created_at"`
	MapUpdatedAt time.Time `gorm:"column:map_updated_at"`
}

type VariableModel struct {
	Name      string
	Metadata  []byte
	Groups    []string
	Behaviour string
	Value     []byte
}

type Model struct {
	Fields          []string
	Values          VariableModel
	MapName         string
	ID              string
	ShortID         string
	VariableName    string
	VariableID      string
	VariableShortID string
	ProjectID       string
	Locale          string
}

func NewModel(projectId, locale, mapName, id, shortID, variableName, variableID, variableShortID string, fields []string, values VariableModel) Model {
	return Model{
		MapName:         mapName,
		ID:              id,
		ShortID:         shortID,
		VariableID:      variableID,
		VariableShortID: variableShortID,
		Locale:          locale,
		Fields:          fields,
		ProjectID:       projectId,
		Values:          values,
		VariableName:    variableName,
	}
}

type LogicResult struct {
	Entry     MapVariableWithMap
	Locale    string
	ProjectID string
}

func (a *Model) Validate() map[string]string {
	v := map[string]interface{}{
		"groups":           a.Values.Groups,
		"mapName":          a.MapName,
		"id":               a.ID,
		"mapIdExists":      nil,
		"variableIdExists": nil,
		"fieldsValid":      a.Fields,
		"variableName":     a.VariableName,
		"variableID":       a.VariableID,
		"behaviour":        a.Values.Behaviour,
		"projectID":        a.ProjectID,
		"locale":           a.Locale,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("mapName", validation.When(a.MapName != "", validation.RuneLength(1, 200))),
			validation.Key("id", validation.When(a.ID != "", validation.RuneLength(26, 26))),
			validation.Key("mapIdExists", validation.By(func(value interface{}) error {
				name := a.MapName
				shortId := a.ShortID
				id := a.ID

				if name == "" && shortId == "" && id == "" {
					return errors.New("At least one of 'id', 'name' or 'shortID' must be supplied in order to identify this map.")
				}
				return nil
			})),
			validation.Key("variableName", validation.When(a.VariableName != "", validation.RuneLength(1, 200))),
			validation.Key("variableID", validation.When(a.VariableID != "", validation.RuneLength(26, 26))),
			validation.Key("variableIdExists", validation.By(func(value interface{}) error {
				name := a.VariableName
				shortId := a.VariableShortID
				id := a.VariableID

				if name == "" && shortId == "" && id == "" {
					return errors.New("At least one of 'id', 'name' or 'shortID' must be supplied in order to identify this variable.")
				}
				return nil
			})),
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

type Variable struct {
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

	Variable Variable `json:"variable"`
}

func newView(logicResult LogicResult) View {
	variable := logicResult.Entry

	return View{
		ID:        logicResult.Entry.MapID,
		Name:      logicResult.Entry.MapName,
		Locale:    logicResult.Locale,
		ProjectID: logicResult.ProjectID,
		CreatedAt: logicResult.Entry.MapCreatedAt,
		UpdatedAt: logicResult.Entry.MapUpdatedAt,
		Variable: Variable{
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
