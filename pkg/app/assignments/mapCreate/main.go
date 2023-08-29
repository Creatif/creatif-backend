package mapCreate

import (
	"creatif/pkg/app/domain/assignments"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/storage"
)

type Main struct {
	model *AssignValueModel
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

func (c Main) Logic() (LogicModel, error) {
	var assignmentNode assignments.MapNode
	if err := storage.GetBy(assignmentNode.TableName(), "map_node_id", c.model.workingMap.ID, &assignmentNode); err != nil {
		return LogicModel{}, appErrors.NewDatabaseError(err).AddError("MapCreate.Logic.Create", nil)
	}

	assignmentMapValue := assignments.NewMapValueNode(assignmentNode.ID, c.model.Value)
	if err := storage.Create(assignmentMapValue.TableName(), &assignmentMapValue, false); err != nil {
		return LogicModel{}, appErrors.NewDatabaseError(err).AddError("MapCreate.Logic.Create", nil)
	}

	return LogicModel{
		m:              c.model.workingMap,
		assignmentNode: assignmentNode,
		valueNode:      assignmentMapValue,
	}, nil
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

func New(model AssignValueModel) pkg.Job[*AssignValueModel, View, LogicModel] {
	return Main{model: &model}
}
