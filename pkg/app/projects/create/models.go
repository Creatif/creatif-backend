package create

import (
	"creatif/pkg/app/domain"
	"time"
)

type CreateProjectModel struct {
	Name string `json:"name"`
}

func NewCreateProjectModel(name string) CreateProjectModel {
	return CreateProjectModel{Name: name}
}

type ProjectView struct {
	Name      string    `json:"name"`
	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
}

func newProjectView(model domain.Project) ProjectView {
	return ProjectView{
		Name:      model.Name,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}
