package create

import (
	"creatif/pkg/app/domain/assignments"
	"creatif/pkg/app/domain/declarations"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/constants"
	"creatif/pkg/lib/sdk"
	"creatif/pkg/lib/storage"
)

type Main struct {
	model GetNodeModel
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

func (c Main) Logic() (NodeWithValueQuery, error) {
	var nodeWithValueQuery NodeWithValueQuery
	var node declarations.Node

	if sdk.IsValidUuid(c.model.ID) {
		if err := storage.Get((&declarations.Node{}).TableName(), c.model.ID, &node, "ID", "Name", "Type", "Behaviour", "Groups", "Metadata"); err != nil {
			return NodeWithValueQuery{}, appErrors.NewDatabaseError(err).AddError("Node.Get.Logic", nil)
		}

		nodeWithValueQuery.ID = node.ID
		nodeWithValueQuery.Name = node.Name
		nodeWithValueQuery.Type = node.Type
		nodeWithValueQuery.Groups = node.Groups
		nodeWithValueQuery.Behaviour = node.Behaviour
		nodeWithValueQuery.Metadata = node.Metadata
		nodeWithValueQuery.CreatedAt = node.CreatedAt
		nodeWithValueQuery.UpdatedAt = node.UpdatedAt

		var textNode assignments.NodeText
		var boolNode assignments.NodeBoolean
		if nodeWithValueQuery.Type == constants.ValueTextType {
			if res := storage.Gorm().Raw(`
				SELECT ant.* FROM assignments.node_text AS ant
				INNER JOIN assignments.nodes AS an ON an.id = ant.assignment_node_id
				INNER JOIN declarations.nodes AS dn ON dn.id = an.declaration_node_id
				WHERE dn.id = ?
			`, nodeWithValueQuery.ID).
				Scan(&textNode); res.Error != nil {
				return NodeWithValueQuery{}, appErrors.NewDatabaseError(res.Error).AddError("Node.Get.Logic", nil)
			}

			nodeWithValueQuery.Value = textNode.Value
		} else if nodeWithValueQuery.Type == constants.ValueBooleanType {
			if res := storage.Gorm().Raw(`
				SELECT ant.value FROM assignments.node_boolean AS ant
				INNER JOIN assignments.nodes AS an ON an.id = ant.assignment_node_id
				INNER JOIN declarations.nodes AS dn ON dn.id = an.declaration_node_id
				WHERE dn.id = ?
			`, nodeWithValueQuery.ID).Scan(&boolNode); res.Error != nil {
				return NodeWithValueQuery{}, appErrors.NewDatabaseError(res.Error).AddError("Node.Get.Logic", nil)
			}

			nodeWithValueQuery.Value = boolNode.Value
		}
	} else {
		if err := storage.GetBy((&declarations.Node{}).TableName(), "name", c.model.ID, &node, "ID", "Name", "Type", "Behaviour", "Groups", "Metadata"); err != nil {
			return NodeWithValueQuery{}, appErrors.NewDatabaseError(err).AddError("Node.Get.Logic", nil)
		}

		nodeWithValueQuery.ID = node.ID
		nodeWithValueQuery.Name = node.Name
		nodeWithValueQuery.Type = node.Type
		nodeWithValueQuery.Groups = node.Groups
		nodeWithValueQuery.Behaviour = node.Behaviour
		nodeWithValueQuery.Metadata = node.Metadata
		nodeWithValueQuery.CreatedAt = node.CreatedAt
		nodeWithValueQuery.UpdatedAt = node.UpdatedAt

		var textNode assignments.NodeText
		var boolNode assignments.NodeBoolean
		if nodeWithValueQuery.Type == constants.ValueTextType {
			if res := storage.Gorm().Raw(`
				SELECT ant.value FROM assignments.node_text AS ant
				INNER JOIN assignments.nodes AS an ON an.id = ant.assignment_node_id
				INNER JOIN declarations.nodes AS dn ON dn.id = an.declaration_node_id
				WHERE dn.id = ?
			`, nodeWithValueQuery.ID).Scan(&textNode); res.Error != nil {
				return NodeWithValueQuery{}, appErrors.NewDatabaseError(res.Error).AddError("Node.Get.Logic", nil)
			}

			nodeWithValueQuery.Value = textNode.Value
		} else if nodeWithValueQuery.Type == constants.ValueBooleanType {
			if res := storage.Gorm().Raw(`
				SELECT ant.value FROM assignments.node_boolean AS ant
				INNER JOIN assignments.nodes AS an ON an.id = ant.assignment_node_id
				INNER JOIN declarations.nodes AS dn ON dn.id = an.declaration_node_id
				WHERE dn.id = ?
			`, nodeWithValueQuery.ID).Scan(&boolNode); res.Error != nil {
				return NodeWithValueQuery{}, appErrors.NewDatabaseError(res.Error).AddError("Node.Get.Logic", nil)
			}

			nodeWithValueQuery.Value = boolNode.Value
		}

	}

	return nodeWithValueQuery, nil
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

func New(model GetNodeModel) pkg.Job[GetNodeModel, View, NodeWithValueQuery] {
	return Main{model: model}
}
