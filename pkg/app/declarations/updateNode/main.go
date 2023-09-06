package updateNode

import (
	"creatif/pkg/app/domain/declarations"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/storage"
)

type Main struct {
	model Model
}

func (c Main) Validate() error {
	if errs := c.model.Validate(); errs != nil {
		return appErrors.NewValidationError(errs)
	}

	return nil
}

func (c Main) Authenticate() error {
	return nil
}

func (c Main) Authorize() error {
	return nil
}

func (c Main) Logic() (declarations.Node, error) {
	var existing declarations.Node
	if err := storage.GetBy((declarations.Node{}).TableName(), "name", c.model.Name, &existing); err != nil {
		return declarations.Node{}, appErrors.NewNotFoundError(err).AddError("updateNode.Logic", nil)
	}

	if c.model.UpdatingName != "" {
		existing.Name = c.model.UpdatingName
	}

	if len(c.model.Metadata) != 0 {
		existing.Metadata = c.model.Metadata
	}

	if len(c.model.Behaviour) != 0 {
		existing.Behaviour = c.model.Behaviour
	}

	if len(c.model.Groups) != 0 {
		existing.Groups = c.model.Groups
	}

	if res := storage.Gorm().Save(&existing); res.Error != nil {
		return declarations.Node{}, appErrors.NewApplicationError(res.Error).AddError("updateNode.Logic", nil)
	}

	return existing, nil
}

func (c Main) Handle() (View, error) {
	if err := c.Validate(); err != nil {
		return View{}, err
	}

	if err := c.Authenticate(); err != nil {
		return View{}, err
	}

	if err := c.Authorize(); err != nil {
		return View{}, err
	}

	model, err := c.Logic()

	if err != nil {
		return View{}, err
	}

	return newView(model), nil
}

func New(model Model) pkg.Job[Model, View, declarations.Node] {
	return Main{model: model}
}
