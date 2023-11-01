package getMap

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/locales"
	"creatif/pkg/lib/sdk"
	"encoding/json"
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"strings"
	"time"
)

var validFields = []string{
	"behaviour",
	"metadata",
	"groups",
	"value",
	"created_at",
	"updated_at",
}

type Model struct {
	Name      string
	ID        string
	ShortID   string
	Fields    []string
	ProjectID string
	Locale    string
	Groups    []string

	validFields []string
}

func NewModel(projectId, locale, name, id, shortId string, fields []string, groups []string) Model {
	return Model{
		Name:        name,
		ID:          id,
		ShortID:     shortId,
		ProjectID:   projectId,
		Locale:      locale,
		Fields:      fields,
		Groups:      groups,
		validFields: validFields,
	}
}

type LogicModel struct {
	variableMap declarations.Map
	variables   []Variable
}

type Variable struct {
	ID string `json:"id" gorm:"primarykey"`

	Name      string         `json:"name" gorm:"index;uniqueIndex:unique_variable"`
	Value     datatypes.JSON `json:"value"`
	Behaviour string         `json:"behaviour"`
	Groups    pq.StringArray `json:"groups" gorm:"type:text[]"`
	Metadata  datatypes.JSON `json:"metadata"`

	CreatedAt time.Time `json:"createdAt" gorm:"<-:createProject"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type View struct {
	ID        string                   `json:"id"`
	Name      string                   `json:"name"`
	ProjectID string                   `json:"projectID"`
	Locale    string                   `json:"locale"`
	Variables []map[string]interface{} `json:"variables"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func newView(model LogicModel, returnFields []string, locale string) View {
	m := make([]map[string]interface{}, 0)

	for _, n := range model.variables {
		o := make(map[string]interface{})

		o["id"] = n.ID
		o["name"] = n.Name

		for _, f := range returnFields {
			if f == "groups" {
				o["groups"] = n.Groups
			}

			if f == "value" {
				o["value"] = n.Value
			}

			if f == "behaviour" {
				o["behaviour"] = n.Behaviour
			}

			if f == "metadata" {
				o["metadata"] = n.Metadata
			}

			if f == "value" {
				o["value"] = n.Value
			}

			if f == "created_at" {
				o["createdAt"] = n.CreatedAt
			}

			if f == "updated_at" {
				o["updatedAt"] = n.UpdatedAt
			}
		}

		m = append(m, o)
	}

	return View{
		ID:        model.variableMap.ID,
		ProjectID: model.variableMap.ProjectID,
		Locale:    locale,
		Name:      model.variableMap.Name,
		Variables: m,
		CreatedAt: model.variableMap.CreatedAt,
		UpdatedAt: model.variableMap.UpdatedAt,
	}
}

func (a *Model) Validate() map[string]string {
	v := map[string]interface{}{
		"name":        a.Name,
		"id":          a.ID,
		"idExists":    nil,
		"projectID":   a.ProjectID,
		"locale":      a.Locale,
		"fieldsValid": a.Fields,
	}

	if err := validation.Validate(v,
		validation.Map(
			// Name cannot be empty, and the length must be between 5 and 20.
			validation.Key("name", validation.When(a.Name != "", validation.RuneLength(1, 200))),
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
			validation.Key("projectID", validation.Required, validation.RuneLength(26, 26)),
			validation.Key("locale", validation.Required, validation.By(func(value interface{}) error {
				t := value.(string)

				if !locales.ExistsByAlpha(t) {
					return errors.New(fmt.Sprintf("Locale '%s' does not exist.", t))
				}

				return nil
			})),
			validation.Key("fieldsValid", validation.By(func(value interface{}) error {
				fields := value.([]string)
				vFields := a.validFields

				if len(fields) > 0 {
					for _, f := range fields {
						if !sdk.Includes(vFields, f) {
							return errors.New(fmt.Sprintf("%s is not a valid field to return. Valid fields are %s", f, strings.Join(a.validFields, ", ")))
						}
					}
				}

				return nil
			})),
		),
	); err != nil {
		var e map[string]string
		b, err := json.Marshal(err)
		if err != nil {
			return map[string]string{
				"unrecoverable": "An internal validation error occurred. This should not happen. Please, submit a bug.",
			}
		}

		if err := json.Unmarshal(b, &e); err != nil {
			return map[string]string{
				"unrecoverable": "An internal validation error occurred. This should not happen. Please, submit a bug.",
			}
		}

		return e
	}

	return nil
}
