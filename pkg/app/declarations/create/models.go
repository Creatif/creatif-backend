package create

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/constants"
	"creatif/pkg/lib/sdk"
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/segmentio/ksuid"
	"time"
)

type ValidationLength struct {
	Min   int `json:"min"`
	Max   int `json:"max"`
	Exact int `json:"exact"`
}

type NodeValidation struct {
	Required    bool
	Length      ValidationLength
	ExactValue  string
	ExactValues []string
	IsDate      bool
}

type CreateNodeModel struct {
	Name       string         `json:"name"`
	Type       string         `json:"type"`
	Metadata   []byte         `json:"metadata"`
	Groups     []string       `json:"groups"`
	Behaviour  string         `json:"behaviour"`
	Validation NodeValidation `json:"validation"`
}

func NewCreateNodeModel(name, behaviour string, groups []string, metadata []byte, validation NodeValidation) CreateNodeModel {
	return CreateNodeModel{
		Name:       name,
		Behaviour:  behaviour,
		Groups:     groups,
		Validation: validation,
		Metadata:   metadata,
	}
}

func (a *CreateNodeModel) Validate() map[string]string {
	v := map[string]interface{}{
		"name":           a.Name,
		"groups":         a.Groups,
		"behaviour":      a.Behaviour,
		"nodeValidation": a.Validation,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("name", validation.Required, validation.RuneLength(1, 200)),
			validation.Key("groups", validation.When(len(a.Groups) != 0, validation.Each(validation.RuneLength(1, 200)))),
			validation.Key("behaviour", validation.Required, validation.By(func(value interface{}) error {
				t := value.(string)

				if t != constants.ReadonlyBehaviour && t != constants.ModifiableBehaviour {
					return errors.New(fmt.Sprintf("Invalid value for behaviour. Node type can be 'modifiable' or 'readonly'"))
				}

				return nil
			})),
			validation.Key("nodeValidation", validation.Required, validation.By(func(value interface{}) error {
				t := value.(NodeValidation)

				if len(t.ExactValue) > 200 {
					return errors.New("Invalid value for validation.exactValue. validation.exactValue cannot have more than 200 characters")
				}

				for _, r := range t.ExactValues {
					if len(r) > 200 {
						return errors.New("Invalid value for validation.exactValues. Every entry in validation.exactValues array cannot be more than 200 characters")
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
	ID        ksuid.KSUID            `json:"id"`
	Name      string                 `json:"name"`
	Groups    []string               `json:"groups"`
	Behaviour string                 `json:"behaviour"`
	Metadata  map[string]interface{} `json:"metadata"`

	CreatedAt time.Time `gorm:"<-:create" json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func newView(model declarations.Node) View {
	return View{
		ID:        model.ID,
		Name:      model.Name,
		Groups:    model.Groups,
		Behaviour: model.Behaviour,
		Metadata:  sdk.UnmarshalToMap([]byte(model.Metadata)),
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}
