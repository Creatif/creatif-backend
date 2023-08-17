package create

import (
	"creatif/pkg/app/assignments/create/services"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
)

type Create struct {
	model *CreateNodeModel
}

func (c Create) Validate() error {
	if errs := c.model.Validate(); errs != nil {
		return appErrors.NewValidationError(errs)
	}

	return nil
}

func (c Create) Authenticate() error {
	return nil
}

func (c Create) Authorize() error {
	return nil
}

func (c Create) Logic() (services.AssignmentCreateResult, error) {
	service := services.NewAssignmentCreate(c.model.declarationNode.Name, c.model.declarationNode.Type, c.model.Value, c.model.declarationNode.ID)
	model, err := service.CreateOrUpdate()

	if err != nil {
		return services.AssignmentCreateResult{}, appErrors.NewDatabaseError(err).AddError("Node.Create.Logic", nil)
	}

	return model, nil
}

func (c Create) Handle() (View, error) {
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

	return newView(c.model.declarationNode, model.Value), nil
}

func New(model *CreateNodeModel) pkg.Job[*CreateNodeModel, View, services.AssignmentCreateResult] {
	return Create{model: model}
}
