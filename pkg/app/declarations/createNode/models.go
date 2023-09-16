package createNode

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/constants"
	"creatif/pkg/lib/sdk"
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"time"
)

type Model struct {
	Name      string   `json:"name"`
	Type      string   `json:"type"`
	Metadata  []byte   `json:"metadata"`
	Groups    []string `json:"groups"`
	Behaviour string   `json:"behaviour"`
	Value     []byte   `json:"value"`
}

func NewModel(name, behaviour string, groups []string, metadata []byte, value []byte) Model {
	return Model{
		Name:      name,
		Behaviour: behaviour,
		Groups:    groups,
		Metadata:  metadata,
		Value:     value,
	}
}

func (a *Model) Validate() map[string]string {
	v := map[string]interface{}{
		"name":      a.Name,
		"groups":    a.Groups,
		"behaviour": a.Behaviour,
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
		),
	); err != nil {
		return sdk.ErrorToResponseError(err)
	}

	return nil
}

type View struct {
	ID        string                 `json:"id"`
	Name      string                 `json:"name"`
	Groups    []string               `json:"groups"`
	Behaviour string                 `json:"behaviour"`
	Metadata  map[string]interface{} `json:"metadata"`
	Value     []byte                 `json:"value"`

	CreatedAt time.Time `gorm:"<-:createNode" json:"createdAt"`
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
