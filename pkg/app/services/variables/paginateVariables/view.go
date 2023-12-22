package paginateVariables

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/locales"
	"creatif/pkg/app/services/variables/paginateVariables/pagination"
	"creatif/pkg/lib/sdk"
	"time"
)

type LogicModel struct {
	variables      []Variable
	paginationInfo pagination.PaginationInfo
}

type View struct {
	ID        string      `json:"id"`
	ProjectID string      `json:"projectID"`
	Locale    string      `json:"locale"`
	Name      string      `json:"name"`
	Groups    []string    `json:"groups"`
	Behaviour string      `json:"behaviour"`
	Metadata  interface{} `json:"metadata"`
	Value     interface{} `json:"value"`

	CreatedAt time.Time `gorm:"<-:createProject" json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func newView(models []declarations.Variable) []View {
	return sdk.Map(models, func(idx int, value declarations.Variable) View {
		locale, _ := locales.GetAlphaWithID(value.LocaleID)
		return View{
			ID:        value.ID,
			ProjectID: value.ProjectID,
			Locale:    locale,
			Name:      value.Name,
			Groups:    value.Groups,
			Value:     value.Value,
			Behaviour: value.Behaviour,
			Metadata:  value.Metadata,
			CreatedAt: value.CreatedAt,
			UpdatedAt: value.UpdatedAt,
		}
	})
}
