package create

import (
	"creatif/pkg/lib/sdk"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"time"
)

type Model struct {
	ProjectID string
	Data      []byte
}

type LogicModel struct {
	ID            string
	Data          interface{}
	CreatedAt     time.Time
	ActivityAdded bool
}

func NewModel(projectId string, data []byte) Model {
	return Model{
		ProjectID: projectId,
		Data:      data,
	}
}

func (a Model) Validate() map[string]string {
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

type View struct {
	ID            string      `json:"structures"`
	Data          interface{} `json:"data"`
	CreatedAt     time.Time   `json:"createdAt"`
	ActivityAdded bool        `json:"activityAdded"`
}

func newView(model LogicModel) View {
	return View{
		ID:            model.ID,
		Data:          model.Data,
		CreatedAt:     model.CreatedAt,
		ActivityAdded: model.ActivityAdded,
	}
}
