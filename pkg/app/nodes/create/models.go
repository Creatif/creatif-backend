package create

import (
	"creatif/pkg/app/domain"
	"time"
)

type CreateNodeModel struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	Group     string `json:"group"`
	Behaviour string `json:"behaviour"`
}

func NewCreateNodeModel(name string) CreateNodeModel {
	return CreateNodeModel{Name: name}
}

type View struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	Group     string `json:"group"`
	Behaviour string `json:"behaviour"`

	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
}

func newView(model domain.Node) View {
	return View{
		Name:      model.Name,
		Type:      model.Type,
		Group:     model.Group,
		Behaviour: model.Behaviour,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}
