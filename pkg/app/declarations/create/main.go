package create

import (
	"creatif/pkg/app/domain/assignments"
	"creatif/pkg/app/domain/declarations"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/storage"
	"gorm.io/gorm"
)

type Create struct {
	model CreateNodeModel
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

func (c Create) Logic() (declarations.Node, error) {
	model := declarations.NewNode(c.model.Name, c.model.Behaviour, c.model.Groups, c.model.Metadata)

	if err := storage.Transaction(func(tx *gorm.DB) error {
		if res := tx.Create(&model); res.Error != nil {
			return res.Error
		}

		assignmentModel := assignments.NewNode(model.Name, model.ID)
		if res := tx.Create(&assignmentModel); res.Error != nil {
			return res.Error
		}

		valueModel := assignments.NewValueNode(assignmentModel.ID, nil)
		if res := tx.Create(&valueModel); res.Error != nil {
			return res.Error
		}

		return nil
	}); err != nil {
		return declarations.Node{}, appErrors.NewDatabaseError(err).AddError("Node.Create.Logic", nil)
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

func New(model CreateNodeModel) pkg.Job[CreateNodeModel, View, declarations.Node] {
	return Create{model: model}
}
