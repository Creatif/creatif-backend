package getMapGroups

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/sdk"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/lib/pq"
)

type LogicModel struct {
	Groups pq.StringArray `gorm:"type:text[];not null"`
}

type Model struct {
	Name      string
	ItemID    string
	ProjectID string
}

func NewModel(name, itemId, projectID string) Model {
	return Model{
		Name:      name,
		ItemID:    itemId,
		ProjectID: projectID,
	}
}

func newView(model []declarations.Group) []string {
	return sdk.Map(model, func(idx int, value declarations.Group) string {
		return value.Name
	})
}

func (a *Model) Validate() map[string]string {
	v := map[string]interface{}{
		"id":        a.Name,
		"itemId":    a.ItemID,
		"projectID": a.ProjectID,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("id", validation.Required),
			validation.Key("itemId", validation.Required),
			validation.Key("projectID", validation.Required, validation.RuneLength(26, 26)),
		),
	); err != nil {
		return sdk.ErrorToResponseError(err)
	}

	return nil
}
