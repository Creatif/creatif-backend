package createVariable

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
	"time"
)

type Model struct {
	Name      string
	Metadata  []byte
	Groups    []string
	Behaviour string
	Value     []byte
	ProjectID string
	Locale    string
}

func NewModel(projectId, locale, name, behaviour string, groups []string, metadata []byte, value []byte) Model {
	return Model{
		Name:      name,
		Locale:    locale,
		ProjectID: projectId,
		Behaviour: behaviour,
		Groups:    groups,
		Metadata:  metadata,
		Value:     value,
	}
}

func (a Model) Validate() map[string]string {
	v := map[string]interface{}{
		"name":      a.Name,
		"groups":    a.Groups,
		"behaviour": a.Behaviour,
		"locale":    a.Locale,
		"projectID": a.ProjectID,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("projectID", validation.Required, validation.RuneLength(26, 26)),
			validation.Key("name", validation.Required, validation.RuneLength(1, 200), validation.By(func(value interface{}) error {
				name := value.(string)

				var variable declarations.Variable
				res := storage.Gorm().Where("name = ? AND project_id = ?", name, a.ProjectID).Select("ID").First(&variable)

				if errors.Is(res.Error, gorm.ErrRecordNotFound) {
					return nil
				}

				if res.Error != nil {
					return errors.New(fmt.Sprintf("Record with name '%s' already exists", name))
				}

				if variable.ID != "" {
					return errors.New(fmt.Sprintf("Record with name '%s' already exists", name))
				}

				return nil
			})),
			validation.Key("groups", validation.When(len(a.Groups) != 0, validation.Each(validation.RuneLength(1, 200))), validation.By(func(value interface{}) error {
				groups := value.([]string)
				if len(groups) > 20 {
					return errors.New(fmt.Sprintf("Invalid number of groups for '%s'. Maximum number of groups per variable is 20.", a.Name))
				}
				for _, g := range a.Groups {
					if len(g) > 100 {
						return errors.New(fmt.Sprintf("Invalid group length for '%s'. Maximum number of characters per groups is 100.", g))
					}
				}

				return nil
			})),
			validation.Key("behaviour", validation.Required, validation.By(func(value interface{}) error {
				t := value.(string)

				if t != constants.ReadonlyBehaviour && t != constants.ModifiableBehaviour {
					return errors.New(fmt.Sprintf("Invalid value for behaviour. Variable type can be 'modifiable' or 'readonly'"))
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
	ID        string      `json:"id"`
	ProjectID string      `json:"projectID"`
	Locale    string      `json:"locale"`
	ShortID   string      `json:"shortID"`
	Name      string      `json:"name"`
	Groups    []string    `json:"groups"`
	Behaviour string      `json:"behaviour"`
	Metadata  interface{} `json:"metadata"`
	Value     interface{} `json:"value"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func newView(model declarations.Variable, locale string) View {
	return View{
		ID:        model.ID,
		ProjectID: model.ProjectID,
		ShortID:   model.ShortID,
		Locale:    locale,
		Name:      model.Name,
		Groups:    model.Groups,
		Metadata:  model.Metadata,
		Value:     model.Value,
		Behaviour: model.Behaviour,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}
