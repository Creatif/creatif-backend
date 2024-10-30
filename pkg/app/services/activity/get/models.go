package get

import (
	"creatif/pkg/lib/sdk"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gorm.io/datatypes"
	"time"
)

type Model struct {
	ProjectID string
}

type LogicModel struct {
	ID        string
	Data      datatypes.JSON `gorm:"type:jsonb"`
	CreatedAt time.Time
}

func NewModel(projectId string) Model {
	return Model{
		ProjectID: projectId,
	}
}

type View struct {
	ID        string      `json:"id"`
	Data      interface{} `json:"data"`
	CreatedAt time.Time   `json:"createdAt"`
}

func newView(model []LogicModel) []View {
	return sdk.Map(model, func(idx int, value LogicModel) View {
		return View{
			ID:        value.ID,
			Data:      value.Data,
			CreatedAt: value.CreatedAt,
		}
	})
}

func (a *Model) Validate() map[string]string {
	v := map[string]interface{}{
		"projectID": a.ProjectID,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("projectID", validation.Required, validation.RuneLength(27, 27)),
		),
	); err != nil {
		return sdk.ErrorToResponseError(err)
	}

	return nil
}
