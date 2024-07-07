package updateMapVariable

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
	"locale",
	"behaviour",
	"references",
	"value",
}

type VariableModel struct {
	Name      string
	Locale    string
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
	References   []shared.UpdateReference
	ImagePaths   []string
}

func NewModel(projectId, mapName, variableName string, fields []string, values VariableModel, reference []shared.UpdateReference, imagePaths []string) Model {
	return Model{
		MapName:      mapName,
		Fields:       fields,
		ProjectID:    projectId,
		Values:       values,
		VariableName: variableName,
		References:   reference,
		ImagePaths:   imagePaths,
	}
}

type LogicResult struct {
	Variable declarations.MapVariable
	Groups   []declarations.Group
}

func (a *Model) Validate() map[string]string {
	v := map[string]interface{}{
		"groups":       a.Values.Groups,
		"mapName":      a.MapName,
		"fieldsValid":  a.Fields,
		"variableName": a.VariableName,
		"behaviour":    a.Values.Behaviour,
		"projectID":    a.ProjectID,
		"locale":       a.Values.Locale,
		"references":   nil,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("mapName", validation.Required),
			validation.Key("variableName", validation.Required),
			validation.Key("projectID", validation.Required, validation.RuneLength(27, 27)),
			validation.Key("behaviour", validation.By(func(value interface{}) error {
				if !sdk.Includes(a.Fields, "behaviour") {
					return nil
				}

				v := value.(string)
				if v != constants.ReadonlyBehaviour && v != constants.ModifiableBehaviour {
					return errors.New(fmt.Sprintf("Invalid value for behaviour in variable '%s'. Variable behaviour can be 'modifiable' or 'readonly'", v))
				}

				return nil
			})),
			validation.Key("fieldsValid", validation.Required, validation.By(func(value interface{}) error {
				t := value.([]string)

				if len(t) == 0 || len(t) > 8 {
					return errors.New(fmt.Sprintf("Invalid updateable fields. Valid updatable fields are %s", strings.Join(validUpdateableFields, ", ")))
				}

				if !sdk.ArrEqual(t, validUpdateableFields) {
					return errors.New(fmt.Sprintf("Invalid updateable fields. Valid updatable fields are %s", strings.Join(validUpdateableFields, ", ")))
				}

				return nil
			})),
			validation.Key("groups", validation.When(len(a.Values.Groups) != 0, validation.Each(validation.RuneLength(1, 100))), validation.By(func(value interface{}) error {
				if !sdk.Includes(a.Fields, "groups") {
					return nil
				}

				if a.Values.Groups != nil && len(a.Values.Groups) > 20 {
					return errors.New("Maximum number of groups is 20.")
				}

				return nil
			})),
			validation.Key("locale", validation.By(func(value interface{}) error {
				if !sdk.Includes(a.Fields, "locale") {
					return nil
				}

				t := value.(string)

				if !locales.ExistsByAlpha(t) {
					return errors.New(fmt.Sprintf("Locale '%s' does not exist.", t))
				}

				return nil
			})),
			validation.Key("references", validation.By(func(value interface{}) error {
				if len(a.References) == 0 {
					return nil
				}

				for _, ref := range a.References {
					if ref.StructureType != "map" && ref.StructureType != "list" && ref.StructureType != "variable" {
						return errors.New(fmt.Sprintf("Invalid reference. StructureType is invalid. %s given for one of the structure types", ref.StructureType))
					}

					if ref.Name == "" {
						return errors.New("Invalid reference. Name cannot be blank.")
					}

					if ref.VariableID == "" {
						return errors.New("Invalid reference. VariableID cannot be blank.")
					}

					if ref.StructureName == "" {
						return errors.New("Invalid reference. StructureName cannot be blank.")
					}
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
	Groups    []GroupView `json:"groups"`
	ShortID   string      `json:"shortId"`
	Metadata  interface{} `json:"metadata"`
	Behaviour string      `json:"behaviour"`
	Value     interface{} `json:"value"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
}

type GroupView struct {
	ID   string `json:"id"`
	Name string `json:"name"`
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
		ID:     model.Variable.ID,
		Locale: locale,
		Groups: sdk.Map(model.Groups, func(idx int, value declarations.Group) GroupView {
			return GroupView{
				ID:   value.ID,
				Name: value.Name,
			}
		}),
		Name:      model.Variable.Name,
		Behaviour: model.Variable.Behaviour,
		Metadata:  m,
		Value:     v,
		CreatedAt: model.Variable.CreatedAt,
		UpdatedAt: model.Variable.UpdatedAt,
	}
}
