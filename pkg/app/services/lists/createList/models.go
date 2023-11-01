package createList

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/locales"
	"creatif/pkg/lib/sdk"
	"creatif/pkg/lib/storage"
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gorm.io/gorm"
	"time"
)

type Variable struct {
	Name      string
	Metadata  []byte
	Groups    []string
	Behaviour string
	Value     []byte
}

type Model struct {
	Name      string
	ProjectID string
	Locale    string
	Variables []Variable
}

func NewModel(projectId, locale, name string, variables []Variable) Model {
	return Model{
		Name:      name,
		Locale:    locale,
		ProjectID: projectId,
		Variables: variables,
	}
}

func (a Model) Validate() map[string]string {
	v := map[string]interface{}{
		"name":        a.Name,
		"projectID":   a.ProjectID,
		"locale":      a.Locale,
		"variableLen": len(a.Variables),
		"groups":      nil,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("projectID", validation.Required, validation.RuneLength(26, 26)),
			validation.Key("name", validation.Required, validation.RuneLength(1, 200), validation.By(func(value interface{}) error {
				name := value.(string)

				var variable declarations.List
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
			validation.Key("variableLen", validation.By(func(value interface{}) error {
				l := value.(int)

				if l > 1000 {
					return errors.New("The number of variables when creating a list cannot be higher than 1000.")
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
			validation.Key("groups", validation.By(func(value interface{}) error {
				for _, variable := range a.Variables {
					if len(variable.Groups) > 20 {
						return errors.New(fmt.Sprintf("Invalid number of groups for '%s'. Maximum number of groups per variable is 20.", variable.Name))
					}

					for _, g := range variable.Groups {
						if len(g) > 100 {
							return errors.New(fmt.Sprintf("Invalid group length for '%s'. Maximum number of characters per groups is 200.", g))
						}
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
	ID        string `json:"id"`
	ShortID   string `json:"shortID"`
	ProjectID string `json:"projectID"`
	Name      string `json:"name"`
	Locale    string `json:"locale"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func newView(model declarations.List, locale string) View {
	return View{
		ID:        model.ID,
		ShortID:   model.ShortID,
		ProjectID: model.ProjectID,
		Locale:    locale,
		Name:      model.Name,

		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}
