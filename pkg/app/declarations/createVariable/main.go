package createVariable

import (
	"creatif/pkg/app/domain/declarations"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/storage"
)

type Create struct {
	model Model
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

func (c Create) Logic() (declarations.Variable, error) {
	model := declarations.NewVariable(c.model.Name, c.model.Behaviour, c.model.Groups, c.model.Metadata, c.model.Value)

	if err := storage.Create((declarations.Variable{}).TableName(), &model, false); err != nil {
		return declarations.Variable{}, appErrors.NewDatabaseError(err).AddError("createVariable.Logic", nil)
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

	return newView(model), nil
}

func New(model Model) pkg.Job[Model, View, declarations.Variable] {
	return Create{model: model}
}
