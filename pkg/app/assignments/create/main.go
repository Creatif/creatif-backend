package create

import (
	"creatif/pkg/app/domain/assignments"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/storage"
	"github.com/google/uuid"
	"gorm.io/gorm"
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

func (c Create) Logic() (assignments.Node, error) {
	model := assignments.NewNode(c.model.Name)
	model.DeclarationNodeID = c.model.declarationNode.ID

	err := storage.Transaction(func(tx *gorm.DB) error {
		var id uuid.UUID
		if c.model.Type == "text" {
			requestModel := c.model.Value.(AssignNodeTextModel)
			textModel := assignments.NewNodeText(requestModel.Value)
			if err := storage.Create(textModel.TableName(), &assignments.NodeText{
				Value: textModel.Value,
			}); err != nil {
				return err
			}

			id, _ = uuid.FromBytes([]byte(textModel.ID))
		} else if c.model.Type == "boolean" {
			requestModel := c.model.Value.(AssignNodeBooleanModel)
			textModel := assignments.NewNodeBoolean(requestModel.Value)
			if err := storage.Create(textModel.TableName(), &assignments.NodeBoolean{
				Value: textModel.Value,
			}); err != nil {
				return err
			}

			id, _ = uuid.FromBytes([]byte(textModel.ID))
		}

		model.ValueID = id
		if err := storage.Create(model.TableName(), &model); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return assignments.Node{}, appErrors.NewDatabaseError(err).AddError("Node.Create.Logic", nil)
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

func New(model *CreateNodeModel) pkg.Job[*CreateNodeModel, View, assignments.Node] {
	return Create{model: model}
}
