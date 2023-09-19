package createVariable

import (
	"creatif/pkg/app/domain/declarations"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/sdk"
	"creatif/pkg/lib/storage"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	var metadata []byte
	var value []byte
	if len(c.model.Metadata) > 0 {
		m, err := sdk.CovertToGeneric(c.model.Metadata)
		if err != nil {
			return declarations.Variable{}, appErrors.NewApplicationError(err)
		}

		metadata = m
	}

	if len(c.model.Value) > 0 {
		m, err := sdk.CovertToGeneric(c.model.Value)
		if err != nil {
			return declarations.Variable{}, appErrors.NewApplicationError(err)
		}

		value = m
	}

	model := declarations.NewVariable(c.model.Name, c.model.Behaviour, c.model.Groups, metadata, value)
	res := storage.Gorm().Model(&model).Clauses(clause.Returning{Columns: []clause.Column{
		{Name: "id"},
		{Name: "name"},
		{Name: "behaviour"},
		{Name: "metadata"},
		{Name: "value"},
		{Name: "groups"},
		{Name: "created_at"},
		{Name: "updated_at"},
	}}).Create(&model)

	if errors.Is(res.Error, gorm.ErrRecordNotFound) || res.RowsAffected == 0 {
		return declarations.Variable{}, appErrors.NewNotFoundError(res.Error).AddError("createVariable.Logic", nil)
	} else if res.Error != nil {
		return declarations.Variable{}, appErrors.NewDatabaseError(res.Error).AddError("createVariable.Logic", nil)
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
