package create

import (
	"creatif/pkg/app/domain"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/storage"
)

type Create struct {
	model CreateProjectModel
}

func (c Create) Validate() error {
	return nil
}

func (c Create) Authenticate() error {
	return nil
}

func (c Create) Authorize() error {
	return nil
}

func (c Create) Logic() (domain.Project, error) {
	model := domain.NewProject(c.model.Name)

	if err := storage.Create(model.WithSchema(), &model); err != nil {
		return domain.Project{}, appErrors.NewDatabaseError(err).AddError("Project.Create.Logic", nil)
	}

	return model, nil
}

func (c Create) Handle() (ProjectView, error) {
	if err := c.Validate(); err != nil {
		return ProjectView{}, err
	}

	if err := c.Authenticate(); err != nil {
		return ProjectView{}, err
	}

	if err := c.Authorize(); err != nil {
		return ProjectView{}, err
	}

	model, err := c.Logic()

	if err != nil {
		return ProjectView{}, err
	}

	return newProjectView(model), nil
}

func New(model CreateProjectModel) pkg.Job[CreateProjectModel, ProjectView, domain.Project] {
	return Create{model: model}
}
